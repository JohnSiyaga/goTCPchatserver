package main

type commandID int

// List of Command IDs
const (
	nickCMD commandID = iota // Auto increment command id for all consts starting from 0
	joinCMD
	roomListCMD
	msgCMD
	privateMsgCMD
	quitCMD
	helpCMD
)

type command struct {
	id     commandID
	client *client
	args   []string
}
