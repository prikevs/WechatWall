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
			log.Info("another verifier comes online,", client)
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				log.Info("verifier,", client, "goes offline")
				delete(h.clients, client)
				close(client.send)
			}
		// case verified := <-h.verified:
		// verified message, handle here
		case <-h.broadcast:
			if !LoadNeedVerification() {
				select {
				case msg := <-h.readymsgs:
					log.Info("you use non-verification mode, send message to vmq directly")
					if err := libredis.PublishClassToMQ(&msg, vMQ); err != nil {
						log.Error("failed to publish to the back of vmq:", err)
						break
					}
				default:
				}
				break
			}
			for client := range h.clients {
				var msg libredis.Msg
				empty := false
				empty_reuse := false
				select {
				case msg = <-reuse:
				default:
					empty_reuse = true
				}

				if empty_reuse {
					select {
					case msg = <-h.readymsgs:
					default:
						empty = true
					}
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
					MsgId:      msg.Key(),
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
					log.Info("msg", msg.MsgId, "sent to verifier")
					if err := libredis.SetClassToMapWithTTL(
						&msg, pMsgsMap, LoadMaxMsgWaitingTime()); err != nil {
						log.Error("failed to set message from %s to waiting map:", err)
					}
				default:
					log.Warning("something happened when writing to verifier",
						client, "make it offline")
					close(client.send)
					delete(h.clients, client)
					// reuse the message next time
					reuse <- msg
				}

			}
		}
	}
}
