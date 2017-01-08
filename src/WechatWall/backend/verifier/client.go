package verifier

import (
	"WechatWall/libredis"

	"github.com/gorilla/websocket"

	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader websocket.Upgrader

func init() {
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	if !StrictOrigin {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	}
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func sendNotification(msg libredis.Msg) {
	content := fmt.Sprintf(NotificationMessage, msg.Content)
	ntf := &libredis.Msg{
		UserOpenid: msg.UserOpenid,
		Content:    content,
	}
	if err := libredis.PublishClassToMQ(ntf, sMQ); err != nil {
		log.Error("failed to publish message to smq,", err)
		return
	}
}

// TODO: use NLP to determine whether to override
func updateLastVerifiedMsg(msg *libredis.Msg) {
	if err := lvmMap.Set(msg.UserOpenid, msg.Content); err != nil {
		log.Warning("failed to save last verified message")
	}
}

func handleVMsg(data []byte) error {
	recvm := &VMsgRecvd{}
	if err := json.Unmarshal(data, recvm); err != nil {
		return err
	}
	if recvm.MsgId == "" {
		return errors.New("invalid parameter msgid")
	}

	msg := &libredis.Msg{}
	if err := libredis.GetClassFromMap(
		recvm.MsgId, msg, pMsgsMap); err != nil {
		return errors.New("failed to get message info, maybe due to TTL timeout: " + err.Error())
	}

	// pass set
	if _, err := passSet.Add(msg.UserOpenid); err != nil {
		log.Warning("failed to add user to passSet")
	}

	// update last verified msg
	updateLastVerifiedMsg(msg)

	msg.VerifiedTime = recvm.VerifiedTime
	// publish to wechat wall (vMQ)
	if recvm.ShowNow == true {
		if err := libredis.PublishRClassToMQ(msg, vMQ); err != nil {
			return errors.New("failed to publish to the front of vmq: " + err.Error())
		}
	} else {
		if err := libredis.PublishClassToMQ(msg, vMQ); err != nil {
			return errors.New("failed to publish to the back of vmq: " + err.Error())
		}
	}
	log.Info("Added message from user", msg.Username, "openid:", msg.UserOpenid, "to VMQ")

	if LoadSendNotification() {
		// send message to smq, notify user that msg sent
		go sendNotification(*msg)
	}
	// delete message from pMsgsMap
	if err := libredis.DelClassFromMap(msg, pMsgsMap); err != nil {
		log.Warning("failed to delete msg from pending msgs map, we can wait for TTL")
	}
	return nil
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Errorf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		log.Debug("received message:", string(message))
		// parse this message and add it to vqueue
		if err := handleVMsg(message); err != nil {
			log.Error("cannot handle message: ", err)
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	go client.writePump()
	client.readPump()
}
