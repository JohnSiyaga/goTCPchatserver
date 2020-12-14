package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	nick     string
	room     *room
	commands chan<- command
}

func (c *client) readInput() {
	// Loop for client messages
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')

		// Exit loop if errors encountered
		if err != nil {
			return
		}

		// Parse input
		msg = strings.Trim(msg, "\r\n")

		// Get command from user message
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		// Inform server of command
		switch cmd {
		case "/nick":
			c.commands <- command{
				id:     nickCMD,
				client: c,
				args:   args,
			}
		case "/join":
			c.commands <- command{
				id:     joinCMD,
				client: c,
				args:   args,
			}
		case "/rooms":
			c.commands <- command{
				id:     roomListCMD,
				client: c,
				args:   args,
			}
		case "/msg":
			c.commands <- command{
				id:     msgCMD,
				client: c,
				args:   args,
			}
		case "/pm":
			c.commands <- command{
				id:     privateMsgCMD,
				client: c,
				args:   args,
			}
		case "/quit":
			c.commands <- command{
				id:     quitCMD,
				client: c,
				args:   args,
			}
		case "/help":
			c.commands <- command{
				id:     helpCMD,
				client: c,
				args:   args,
			}
		default: // If command doesn't exist
			c.err(fmt.Errorf("Unknown command: %s", cmd))
			c.msg("Type /help for a full command list \n")
		}
	}
}

// Print errors
func (c *client) err(err error) {
	c.conn.Write([]byte("Error: " + err.Error() + "\n"))
}

// Print message
func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}
