//go:generate msgp
package message

import (
	"fmt"
	"strings"
	"time"
)

// UDPMessage is a binary serializable message for UDP
type UDPMessage struct {
	Context        string      `msg:"context"`
	Name           EventName   `msg:"name"`
	Type           MessageType `msg:"type"`
	Time           time.Time   `msg:"time"`
	Host           string      `msg:"host"`
	FillHostWithIP bool        `msg:"autoIP"`
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
	// MessageNewContextInstance Start a new context instance
	MessageNewTimeline
)

type MessageParseError string

func (m MessageParseError) Error() string {
	return "Message \"" + string(m) + "\" could not be parsed"
}

func ParseMessageType(s string) (MessageType, error) {
	switch s {
	case "MessageStartEvent":
		return MessageStartEvent, nil
	case "MessageEndEvent":
		return MessageEndEvent, nil
	case "MessageSingleEvent":
		return MessageSingleEvent, nil
	default:
		return 0, MessageParseError(s)
	}
}

func (m MessageType) String() string {
	switch m {
	case MessageStartEvent:
		return "MessageStartEvent"
	case MessageEndEvent:
		return "MessageEndEvent"
	case MessageSingleEvent:
		return "MessageSingleEvent"
	case MessageNewTimeline:
		return "MessageNewTimeline"
	default:
		return fmt.Sprintf("Unknown MessageType! <%v>", int(m))
	}
}

// EventName is a binary serialized list of UTF8 strings separated by '.'
// It includes a 2 byte length at the start
type EventName []string

func (e EventName) String() string {
	return strings.Join([]string(e), ".")
}
