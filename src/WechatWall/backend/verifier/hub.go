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

	reuse chan libredis.Msg
}

func newHub(msg <-chan libredis.Msg, bc <-chan bool) *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		broadcast:  bc,
		readymsgs:  msg,
		reuse:      make(chan libredis.Msg, 2),
	}
}

func (h *Hub) handleBroadcast() {
	if !LoadNeedVerification() {
		select {
		case msg := <-h.readymsgs:
			log.Info("you use non-verification mode, send message to vmq directly")
			if err := libredis.PublishClassToMQ(&msg, vMQ); err != nil {
				log.Error("failed to publish to the back of vmq:", err)
			}
		default:
		}
		return
	}
	for client := range h.clients {
		var msg libredis.Msg
		empty := false
		empty_reuse := false
		select {
		case msg = <-h.reuse:
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
			TTL:        int64(LoadMaxMsgWaitingTime().Seconds() * 1000),
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
			h.reuse <- msg
		}
	}
}

func (h *Hub) handleRegister(client *Client) {
	log.Info("another verifier comes online", client, "send TTL messages")
	h.clients[client] = true
	// TODO: send all pending messages
	keys, err := pMsgsMap.Keys()
	if err != nil {
		log.Warning("failed to get keys of pending messages map")
		return
	}
	vmsgs := make([]*VMsgSent, 0)
	for _, key := range keys {
		msg := &libredis.Msg{}
		if err := libredis.GetClassFromMap(key, msg, pMsgsMap); err != nil {
			log.Warning("failed to load waiting message,", key)
			continue
		}
		ttl, err := pMsgsMap.TTL(key)
		if err != nil {
			ttl = LoadMaxMsgWaitingTime()
			log.Warning("cannot get TTL for message", msg.Key(), "use default")
		}
		vmsg := &VMsgSent{
			Username:   msg.Username,
			Openid:     msg.UserOpenid,
			MsgId:      msg.Key(),
			MsgType:    msg.MsgType,
			CreateTime: msg.CreateTime,
			Content:    msg.Content,
			ImgUrl:     utils.BuildImagePath(msg.UserOpenid),
			TTL:        int64(ttl.Seconds() * 1000),
		}
		vmsgs = append(vmsgs, vmsg)
	}
	for _, vmsg := range vmsgs {
		data, err := json.Marshal(vmsg)
		if err != nil {
			log.Warningf("message from %s failed to encode to json",
				vmsg.Openid)
			continue
		}
		select {
		case client.send <- []byte(data):
			log.Info("msg", vmsg.MsgId, "sent to new registered verifier")
		default:
			log.Warning("something happened when writing to verifier",
				client, "make it offline")
			close(client.send)
			delete(h.clients, client)
		}
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.handleRegister(client)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				log.Info("verifier,", client, "goes offline")
				delete(h.clients, client)
				close(client.send)
			}
		// case verified := <-h.verified:
		// verified message, handle here
		case <-h.broadcast:
			h.handleBroadcast()
		}
	}
}
