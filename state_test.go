package main

import (
	"testing"

	"github.com/Redundancy/upto/message"
	"github.com/Redundancy/upto/store"
)

type TestNullStore struct {
}

func (t *TestNullStore) CreateContext(name string) {

}

func (t *TestNullStore) ContextExists(name string) bool {
	return false
}

func (t *TestNullStore) ListContexts() []string {
	return nil
}

func (t *TestNullStore) CreateNewTimeline(contextName string) {
}

func (t *TestNullStore) GetLatestTimeline(contextName string) string {
	return ""
}

func (t *TestNullStore) ListContextTimelines(contextName string) []string {
	return nil
}

func (t *TestNullStore) GetTimelineEvents(contextName, timeline string) []store.TimelineEvent {
	return nil
}

func (t *TestNullStore) AddTimelineEvent(contextName, timeline string, event store.TimelineEvent) {
}

type TestAddEventStore struct {
	Callback func(contextName, timeline string, event store.TimelineEvent)
	*TestNullStore
}

func (s *TestAddEventStore) AddTimelineEvent(contextName, timeline string, event store.TimelineEvent) {
	s.Callback(contextName, timeline, event)
}

func TestMessageHandlerAddsEventAfterStartEnd(t *testing.T) {
	mh := &MessageHandler{}
	calls := 0

	s := &TestAddEventStore{
		Callback: func(contextName, timeline string, event store.TimelineEvent) {
			calls += 1
		},
	}

	startMessage := &message.UDPMessage{
		Context: "ctx",
		Name:    "name",
		Type:    message.MessageStartEvent,
	}

	endMessage := &message.UDPMessage{
		Context: "ctx",
		Name:    "name",
		Type:    message.MessageEndEvent,
	}

	if e := mh.handleMessage(s, startMessage); e != nil {
		t.Error(e)
	}

	if e := mh.handleMessage(s, endMessage); e != nil {
		t.Error(e)
	}

	if calls != 1 {
		t.Fail()
	}
}

func TestMessageHandlerAddsEventAfterStartEndForSameHost(t *testing.T) {
	mh := &MessageHandler{}
	calls := 0

	s := &TestAddEventStore{
		Callback: func(contextName, timeline string, event store.TimelineEvent) {
			calls += 1
		},
	}

	startMessage := &message.UDPMessage{
		Context: "ctx",
		Name:    "name",
		Host:    "a",
		Type:    message.MessageStartEvent,
	}

	endMessage := &message.UDPMessage{
		Context: "ctx",
		Name:    "name",
		Host:    "a",
		Type:    message.MessageEndEvent,
	}

	if e := mh.handleMessage(s, startMessage); e != nil {
		t.Error(e)
	}

	if e := mh.handleMessage(s, endMessage); e != nil {
		t.Error(e)
	}

	if calls != 1 {
		t.Fail()
	}
}

func TestMessageHandlerErrorsAfterStartEndForDifferentHosts(t *testing.T) {
	mh := &MessageHandler{}
	calls := 0

	s := &TestAddEventStore{
		Callback: func(contextName, timeline string, event store.TimelineEvent) {
			calls += 1
		},
	}

	startMessage := &message.UDPMessage{
		Context: "ctx",
		Name:    "name",
		Host:    "a",
		Type:    message.MessageStartEvent,
	}

	endMessage := &message.UDPMessage{
		Context: "ctx",
		Name:    "name",
		Host:    "b",
		Type:    message.MessageEndEvent,
	}

	if e := mh.handleMessage(s, startMessage); e != nil {
		t.Error(e)
	}

	if e := mh.handleMessage(s, endMessage); e == nil {
		t.Error("Should have errored")
	}

	if calls != 0 {
		t.Fail()
	}
}

func TestMessageHandlerErrorsOnUnstartedEvent(t *testing.T) {
	mh := &MessageHandler{}
	s := &TestNullStore{}

	endMessage := &message.UDPMessage{
		Context: "ctx",
		Name:    "name2",
		Type:    message.MessageEndEvent,
	}

	if e := mh.handleMessage(s, endMessage); e == nil {
		t.Fatal("Should have errored")
	}
}

func TestMessageHandlerErrorsOnDoublyStartedEvent(t *testing.T) {
	mh := &MessageHandler{}
	s := &TestNullStore{}

	endMessage := &message.UDPMessage{
		Context: "ctx",
		Name:    "name2",
		Type:    message.MessageStartEvent,
	}

	if e := mh.handleMessage(s, endMessage); e != nil {
		t.Error(e)
	}

	if e := mh.handleMessage(s, endMessage); e == nil {
		t.Fatal("Should have errored")
	}
}
