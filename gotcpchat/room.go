package main

import "net"

type room struct {
	name    string
	members map[net.Addr]*client
}

func (r *room) broadcast(sender *client, msg string) {
	for addr, m := range r.members {
		if addr != sender.conn.RemoteAddr() {
			m.msg(msg)
		}
	}
}

func (r *room) privateMsg(sender *client, recipient string, msg string) {
	for addr, m := range r.members {
		// Don't send message to anonymous people
		if m.nick == "anon" {
			return
		}

		if m.nick == recipient {
			if addr != sender.conn.RemoteAddr() {
				m.msg(msg)
			}
		}
	}
}
