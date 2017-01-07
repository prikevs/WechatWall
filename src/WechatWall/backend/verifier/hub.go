package verifier

import (
	"WechatWall/backend/utils"
	"WechatWall/libredis"

	"encoding/json"
)

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast <-chan bool

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// ready messages
	readymsgs <-chan libredis.Msg
}

func newHub(msg <-chan libredis.Msg, bc <-chan bool) *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		broadcast:  bc,
		readymsgs:  msg,
	}
}

func (h *Hub) run() {
	reuse := make(chan libredis.Msg, 2)
	for {
		select {
		case client := <-h.register:
			log.Info("another client comes online,", client)
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		// case verified := <-h.verified:
		// verified message, handle here
		case <-h.broadcast:
			for client := range h.clients {
				var msg libredis.Msg
				var empty bool
				select {
				case msg = <-reuse:
				case msg = <-h.readymsgs:
				default:
					empty = true
					break
				}

				// no message got, break
				if empty {
					break
				}
				// got message, add message to pending msgs map, set TTL
				log.Info("prepare to send msg", msg.MsgId, "for verification")
				vmsg := &VMsgSent{
					Username:   msg.Username,
					Openid:     msg.UserOpenid,
					MsgId:      msg.MsgId,
					MsgType:    msg.MsgType,
					CreateTime: msg.CreateTime,
					Content:    msg.Content,
					ImgUrl:     utils.BuildImagePath(msg.UserOpenid),
				}
				data, err := json.Marshal(vmsg)
				if err != nil {
					log.Warningf("message from %s failed to encode to json",
						msg.UserOpenid)
					continue
				}

				// prepare to send message
				select {
				case client.send <- []byte(data):
					log.Info("msg", msg.MsgId, "sent to", msg.UserOpenid)
					if err := libredis.SetClassToMapWithTTL(&msg, pMsgsMap, MaxMsgWaitingTime); err != nil {
						log.Error("failed to set message from %s to waiting map:", err)
					}
				default:
					close(client.send)
					delete(h.clients, client)
					// reuse the message next time
					reuse <- msg
				}

			}
		}
	}
}
