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

	reuse chan libredis.Msg
}

func newHub(wallmsgs <-chan libredis.Msg, bc chan bool) *Hub {
	return &Hub{
		broadcast:  bc,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		wallmsgs:   wallmsgs,
		reuse:      make(chan libredis.Msg, 2),
	}
}

func (h *Hub) handleBroadcast() {

}

func (h *Hub) run() {
	reuse := make(chan libredis.Msg, 2)
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Info("another wall comes online,", client)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				log.Info("wall,", client, "goes offline")
				delete(h.clients, client)
				close(client.send)
			}
		case <-h.broadcast:
			var msg libredis.Msg
			if LoadReplay() {
				//TODO
				log.Info("in REPLAY mode, select message from ow set randomly")
				msgid, err := owSet.RandMember()
				if err != nil {
					log.Error("FAILED to get rand member from OW SET!:", err)
					break
				}
				if err := libredis.GetClassFromMap(msgid, &msg, owMap); err != nil {
					log.Error("FAILED to get msg from WATING MSGS MAP!:", err)
					break
				}
			} else {
				empty := false
				empty_reuse := false

				select {
				case msg = <-reuse:
				default:
					empty_reuse = true
				}

				if empty_reuse {
					select {
					case msg = <-h.wallmsgs:
					default:
						empty = true
					}
				}

				// no message got, break
				if empty {
					break
				}

				if LoadReliableMsg() && len(h.clients) == 0 {
					log.Info("there is no wall now, reuse message.")
					select {
					case reuse <- msg:
					default:
						log.Warning("reuse is full, not supposed to be here.")
					}
					break
				}
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
				break
			}
			if !LoadReplay() {
				log.Debug("in Non-Replay mode, save message to owSet and owMap")
				if err := libredis.SetClassToMap(&msg, owMap); err != nil {
					log.Warning("failed to add on wall message to on wall map")
				}
				if _, err := owSet.Add(msg.Key()); err != nil {
					log.Warning("failed to add on wall message to on wall set")
				}
			}
			// prepare to send message
			for client := range h.clients {
				select {
				case client.send <- data:
				default:
					log.Warning("something happened when writing to wall", client, "make it offline")
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
