package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

// Initialize server
type server struct {
	rooms    map[string]*room
	commands chan command
}

// Create new server
func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

// Read and execute incoming commands
func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case nickCMD:
			s.nick(cmd.client, cmd.args)
		case joinCMD:
			s.join(cmd.client, cmd.args)
		case roomListCMD:
			s.listRooms(cmd.client)
		case msgCMD:
			s.msg(cmd.client, cmd.args)
		case privateMsgCMD:
			s.pm(cmd.client, cmd.args)
		case quitCMD:
			s.quit(cmd.client, cmd.args)
		case helpCMD:
			s.help(cmd.client)
		}
	}
}

// Initialize new client connections
func (s *server) newClient(conn net.Conn) {
	log.Printf("New Client has connected: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		nick:     "anon",
		commands: s.commands,
	}

	c.readInput()
}

// Command definitions
func (s *server) nick(c *client, args []string) {
	var previousName string
	previousName = c.nick
	c.nick = args[1] // Get nickname from '/nick <name>'
	c.msg(fmt.Sprintf("New nickname set: %s", c.nick))
	if c.room != nil {
		c.room.broadcast(c, previousName+" has changed their nickname to: "+c.nick)
	}
}

func (s *server) join(c *client, args []string) {
	roomName := args[1] // Get room name from '/join <roomName>'

	// Check if room exists, if not make one
	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}

	// Add client as member of room
	r.members[c.conn.RemoteAddr()] = c

	// Set current room of client
	s.quitCurrentRoom(c) // Quit previous room if switching rooms
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s has joined the room", c.nick))
	c.msg(fmt.Sprintf("Welcome to %s", r.name))
	// Print number of users in chatroom
	c.msg(fmt.Sprintf("There are now currently %v users (including you) in this chatroom", len(r.members)))
}

func (s *server) listRooms(c *client) {
	// Get all rooms
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	// Print rooms to client
	c.msg(fmt.Sprintf("Available rooms: %s", strings.Join(rooms, ", ")))
}

func (s *server) msg(c *client, args []string) {
	// Check if user is currently in a room
	if c.room == nil {
		c.err(errors.New("You must join a room before sending a message"))
		return
	}

	// Broadcast message to room. Take all arguments after '/msg', delimiting by spaces
	c.room.broadcast(c, c.nick+": "+strings.Join(args[1:len(args)], " "))
}

func (s *server) pm(c *client, args []string) {
	// Check if user is currently in a room
	if c.room == nil {
		c.err(errors.New("You must join a room before sending a private message"))
		return
	}

	// Send message to nicknamed connection
	c.room.privateMsg(c, args[1], "Private message from "+c.nick+": "+strings.Join(args[2:len(args)], " "))
}

func (s *server) quit(c *client, args []string) {
	// Quit server and close connection
	log.Printf("Client has disconnected: %s", c.conn.RemoteAddr().String())
	s.quitCurrentRoom(c)
	c.msg("Goodbye!")
	c.conn.Close()
}

func (s *server) help(c *client) {
	// Print command list
	c.msg("/nick <name> - Set your name (Default: anon) \n")
	c.msg("/join <name> - Join a chatroom. If there isn't one, make one \n")
	c.msg("/rooms - Show available rooms \n")
	c.msg("/msg <msg> - Message all people in room \n")
	c.msg("/pm <nick> <msg> - Message specific person in room (if nicknamed) \n")
	c.msg("/quit - Leave chatserver \n")
	c.msg("/help - Show this command list")
}

// Quit room, but not chatserver
func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}
