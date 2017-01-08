package wall

import (
	"WechatWall/backend/utils"
	"WechatWall/libredis"

	"encoding/json"
)

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan bool

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	wallmsgs <-chan libredis.Msg
}

func newHub(wallmsgs <-chan libredis.Msg, bc chan bool) *Hub {
	return &Hub{
		broadcast:  bc,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		wallmsgs:   wallmsgs,
	}
}

func (h *Hub) run() {
	reuse := make(chan libredis.Msg, 2)
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case <-h.broadcast:
			var msg libredis.Msg
			empty := false
			select {
			case msg = <-reuse:
				break
			case msg = <-h.wallmsgs:
			default:
				empty = true
			}

			// no message got, break
			if empty {
				break
			}

			if ReliableMsg && len(h.clients) == 0 {
				log.Info("there is no wall now, reuse message.")
				reuse <- msg
				break
			}

			log.Info("prepare to send msg", msg.MsgId, msg.Content, "to Wall")
			wmsg := &WallMsg{
				MsgId:      msg.MsgId,
				Username:   msg.Username,
				Openid:     msg.UserOpenid,
				MsgType:    msg.MsgType,
				CreateTime: msg.CreateTime,
				Content:    msg.Content,
				ImgUrl:     utils.BuildImagePath(msg.UserOpenid),
			}

			data, err := json.Marshal(wmsg)
			if err != nil {
				log.Warning("message from %s failed to encode to json",
					msg.UserOpenid)
				continue
			}
			// prepare to send message
			for client := range h.clients {
				select {
				case client.send <- data:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
