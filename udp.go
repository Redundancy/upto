package main

import "fmt"

// UDPMessage is a binary serializable message for UDP
type UDPMessage struct {
	Name EventName
	Type MessageType
	Time int64
}

// MessageType is used to indicate the difference between a start and end event
type MessageType int

const (
	// MessageStartEvent should be used as part of UDPMessage when an event is starting
	MessageStartEvent = MessageType(iota)
	// MessageEndEvent should be used as part of UDPMessage when an event is finishing
	MessageEndEvent
	// MessageSingleEvent should be used to indicate a point-in-time event
	MessageSingleEvent
)

func (m MessageType) String() string {
	switch m {
	case MessageStartEvent:
		return "MessageStartEvent"
	case MessageEndEvent:
		return "MessageEndEvent"
	case MessageSingleEvent:
		return "MessageSingleEvent"
	default:
		return fmt.Sprintf("Unknown MessageType! <%v>", int(m))
	}
}

// EventName is a binary serialized list of UTF8 strings separated by '.'
// It includes a 2 byte length at the start
type EventName []string

// MarshalBinary implements encoding.BinaryMarshaler
func (e EventName) MarshalBinary() (data []byte, err error) {
	return
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (e EventName) UnmarshalBinary(data []byte) error {

	return nil
}

func udpMessage(message []byte) {
	fmt.Print("Message:", string(message))
}