package main

import (
	"fmt"
	"sync"

	"github.com/Redundancy/upto/message"
	"github.com/Redundancy/upto/store"
)

// BinaryMessageHandler reads MessagePack encoded messages
type BinaryMessageHandler struct {
	MessageHandler
}

func (h *BinaryMessageHandler) ReceiveMessage(store store.UptoDataStore, ip string, buffer []byte) error {
	m := &message.UDPMessage{}
	m.UnmarshalMsg(buffer)

	if m.FillHostWithIP {
		m.Host = ip
	}

	return h.handleMessage(store, m)
}

type HostMap map[string]store.TimelineEvent
type EventHostMap map[string]HostMap
type ContextEventHostMap map[string]EventHostMap

func (h ContextEventHostMap) getEvent(context, event, host string) (evt store.TimelineEvent, exists bool) {
	c, contextExists := h[context]

	if !contextExists {
		exists = false
		return
	}

	e, eventNameExists := c[event]

	if !eventNameExists {
		exists = false
		return
	}

	evt, exists = e[host]
	return
}

func (h ContextEventHostMap) addEvent(context string, event store.TimelineEvent) error {
	name := event.EventName
	c, contextExists := h[context]

	if !contextExists {
		c = make(EventHostMap)
		h[context] = c
	}

	e, eventNameExists := c[name]

	if !eventNameExists {
		e = make(HostMap)
		c[name] = e
	}

	_, exists := e[event.Host]
	if exists {
		return EventInProgressError(name)
	}

	e[event.Host] = event
	return nil
}

// MessageHandler handles basic message correlation and handling
type MessageHandler struct {
	sync.Mutex
	ContextEventHostMap
}

func (h *MessageHandler) ensureMessagesMapExists() {
	if h.ContextEventHostMap == nil {
		h.ContextEventHostMap = make(ContextEventHostMap)
	}
}

// EventInProgressError is returned when events are already in progress and
// have not completed
type EventInProgressError string

func (e EventInProgressError) Error() string {
	return "Event " + string(e) + " was already in progress"
}

// EventNotStartedError is returned an end event is received for an event that
// has not begun
type EventNotStartedError string

func (e EventNotStartedError) Error() string {
	return "Event " + string(e) + " had not had a start event"
}

func (h *MessageHandler) handleMessage(data store.UptoDataStore, m *message.UDPMessage) error {
	h.Lock()
	defer h.Unlock()

	h.ensureMessagesMapExists()
	name := m.Name

	switch m.Type {
	case message.MessageStartEvent:
		return h.addEvent(m.Context, store.TimelineEvent{
			EventName: name,
			Host:      m.Host,
			Start:     m.Time,
			End:       m.Time,
		})

	case message.MessageEndEvent:
		evt, exists := h.getEvent(m.Context, name, m.Host)
		if !exists {
			return EventNotStartedError(name)
		}

		data.AddTimelineEvent(
			m.Context,
			data.GetLatestTimeline(m.Context),
			store.TimelineEvent{
				EventName: name,
				Host:      m.Host,
				Start:     evt.Start,
				End:       m.Time,
			},
		)

	case message.MessageSingleEvent:
		return fmt.Errorf(
			"Single Events not implemented: %v ",
			name,
		)

	case message.MessageNewTimeline:
		if !data.ContextExists(m.Context) {
			data.CreateContext(m.Context)
		}
		data.CreateNewTimeline(m.Context)
	}

	return nil
}
