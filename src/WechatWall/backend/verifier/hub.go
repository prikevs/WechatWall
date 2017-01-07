package verifier

import (
	"WechatWall/libredis"
)

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan bool

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// ready messages
	readymsgs <-chan libredis.Msg
}

func newHub(ch <-chan libredis.Msg) *Hub {
	return &Hub{
		broadcast:  make(chan bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		readymsgs:  ch,
	}
}

// TODO: pre handle ready messages

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
		// case verified := <-h.verified:
		// verified message, handle here
		case <-h.broadcast:
			for client := range h.clients {
				var msg libredis.Msg
				select {
				case msg = <-reuse:
				case msg = <-h.readymsgs:
				default:
					// no message got, break
					break
				}

				// got message, add message to pending msgs map, set TTL
				data, err := msg.Json()
				if err != nil {
					log.Warningf("message from %s failed to encode to json",
						msg.UserOpenid)
					continue
				}
				select {
				case client.send <- []byte(data):
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
