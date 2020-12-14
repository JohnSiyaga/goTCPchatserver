package main

type commandID int

// List of Command IDs
const (
	nickCMD       commandID = iota // Auto increment command id for all consts starting from 0
	joinCMD                        // 1
	roomListCMD                    // 2
	msgCMD                         // 3
	privateMsgCMD                  // 4
	quitCMD                        // 5
)

type command struct {
	id     commandID
	client *client
	args   []string
}
