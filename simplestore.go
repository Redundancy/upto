package main

import (
	"crypto/rand"
	"fmt"
	"sync"

	"github.com/Redundancy/upto/store"
)

// SimpleMemoryStore keeps all events in memory
// it's intended to be the shortest path to get upto testable
type SimpleMemoryStore struct {
	sync.Mutex
	ContextTimelineEventHostMap
	latestTimeline map[string]string
}

type TimelineMap map[string]EventHostMap
type ContextTimelineEventHostMap map[string]TimelineMap

func (h ContextTimelineEventHostMap) addEvent(context string, timeline string, event store.TimelineEvent) error {
	name := event.EventName
	c, contextExists := h[context]

	if !contextExists {
		c = make(TimelineMap)
		h[context] = c
	}

	t, timelineExists := c[timeline]

	if !timelineExists {
		t = make(EventHostMap)
		c[timeline] = t
	}

	e, eventNameExists := t[name]

	if !eventNameExists {
		e = make(HostMap)
		t[name] = e
	}

	_, exists := e[event.Host]
	if exists {
		return EventInProgressError(name)
	}

	e[event.Host] = event
	return nil
}

func (s *SimpleMemoryStore) CreateContext(name string) {
	s.Lock()
	defer s.Unlock()

	if s.ContextTimelineEventHostMap == nil {
		s.ContextTimelineEventHostMap = make(ContextTimelineEventHostMap)
	}

	s.ContextTimelineEventHostMap[name] = make(TimelineMap)
}

func (s *SimpleMemoryStore) ContextExists(name string) bool {
	s.Lock()
	defer s.Unlock()

	if s.ContextTimelineEventHostMap == nil {
		return false
	}

	_, found := s.ContextTimelineEventHostMap[name]
	return found
}

func (s *SimpleMemoryStore) ListContexts() []string {
	s.Lock()
	defer s.Unlock()

	r := make([]string, 0, len(s.ContextTimelineEventHostMap))

	for context, _ := range s.ContextTimelineEventHostMap {
		r = append(r, context)
	}

	return r
}

func uuid() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func (s *SimpleMemoryStore) CreateNewTimeline(contextName string) {
	s.Lock()
	defer s.Unlock()

	if s.latestTimeline == nil {
		s.latestTimeline = make(map[string]string)
	}

	// TODO: make more than one timeline
	s.latestTimeline[contextName] = uuid()
}

func (s *SimpleMemoryStore) GetLatestTimeline(contextName string) string {
	s.Lock()
	defer s.Unlock()

	if s.latestTimeline == nil {
		return ""
	}

	return s.latestTimeline[contextName]
}

func (s *SimpleMemoryStore) ListContextTimelines(contextName string) []string {
	s.Lock()
	defer s.Unlock()

	if s.ContextTimelineEventHostMap == nil {
		return nil
	}

	m, found := s.ContextTimelineEventHostMap[contextName]

	if !found {
		return nil
	}

	r := make([]string, 0, len(m))

	for k, _ := range m {
		r = append(r, k)
	}

	return r
}

func (s *SimpleMemoryStore) GetTimelineEvents(contextName, timeline string) []store.TimelineEvent {
	s.Lock()
	defer s.Unlock()

	if s.ContextTimelineEventHostMap == nil {
		return nil
	}

	m, found := s.ContextTimelineEventHostMap[contextName]

	if !found {
		return nil
	}

	t, found := m[timeline]

	if !found {
		return nil
	}

	r := make([]store.TimelineEvent, 0)

	for _, hosts := range t {
		for _, event := range hosts {
			r = append(r, event)
		}
	}

	return r
}

func (s *SimpleMemoryStore) AddTimelineEvent(contextName, timeline string, event store.TimelineEvent) {
	s.Lock()
	defer s.Unlock()
	s.ContextTimelineEventHostMap.addEvent(contextName, timeline, event)
}
