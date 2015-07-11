package store

import "time"

type TimelineEvent struct {
	EventName string
	Host      string
	Start     time.Time
	End       time.Time
}

/* UptoDataStore should be implemented to provide a persistence layer
   for upto. Caching can also be provided by the implementor. */
type UptoDataStore interface {
	CreateContext(name string)
	ContextExists(name string) bool
	ListContexts() []string

	CreateNewTimeline(contextName string)
	GetLatestTimeline(contextName string) string
	ListContextTimelines(contextName string) []string

	GetTimelineEvents(contextName, timeline string) []TimelineEvent
	AddTimelineEvent(contextName, timeline string, event TimelineEvent)
}
