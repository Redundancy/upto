package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
)

const (
	currentFolderPrefix = ""
)

// TODO: put a shim REST API here to implement a
func init() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Get("/", rootHandler),
		rest.Get("/contexts", getContexts),
		rest.Get("/contexts/:context", getContext),
		rest.Get("/contexts/:context/:timeline", getTimeline),
	)

	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	mux.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
}

func rootHandler(w rest.ResponseWriter, req *rest.Request) {
	rw := halResponseWriter{w}
	re := MakeHalResponse()
	re.Links["contexts"] = SimpleHREF("contexts")

	rw.WriteJSONResponse(req.URL, re)
}

func getContexts(w rest.ResponseWriter, req *rest.Request) {
	rw := halResponseWriter{w}
	re := MakeHalResponse()

	items := []SimpleHREF{
		currentFolderPrefix + "test",
	}

	for _, context := range datastore.ListContexts() {
		items = append(
			items,
			SimpleHREF(context),
		)
	}

	re.Links["items"] = items
	rw.WriteJSONResponse(req.URL, re)
}

type ContextResponse struct {
	*HalResponse
	Name        string
	Description string
}

func getContext(w rest.ResponseWriter, req *rest.Request) {
	rw := halResponseWriter{w}
	re := &ContextResponse{HalResponse: MakeHalResponse()}

	re.Name = req.PathParam("context")

	timelines := make([]SimpleHREF, 0, len(testTimelines))

	if req.PathParam("context") == "test" {
		for n, _ := range testTimelines {
			timelines = append(
				timelines,
				SimpleHREF(currentFolderPrefix+n),
			)

			re.Links["latest"] = SimpleHREF(currentFolderPrefix + n)
		}
	} else {
		tls := datastore.ListContextTimelines(re.Name)
		for _, timeline := range tls {
			timelines = append(
				timelines,
				SimpleHREF(currentFolderPrefix+timeline),
			)
		}

		latest := datastore.GetLatestTimeline(re.Name)
		if latest != "" {
			re.Links["latest"] = SimpleHREF(currentFolderPrefix + latest)
		}
	}

	re.Links["items"] = timelines
	rw.WriteJSONResponse(req.URL, re)
}

type Event struct {
	Name  string
	Host  string `json:",omitempty"`
	Start time.Time
	End   time.Time
}

var testNow = time.Now()

var testTimelines = map[string][]Event{
	"1": []Event{
		{
			Name:  "File Copy",
			Host:  "10.1.1.345",
			Start: testNow,
			End:   testNow.Add(time.Minute),
		},
		{
			Name:  "DB Update.Proc1",
			Start: testNow,
			End:   testNow.Add(time.Minute * 5),
		},
		{
			Name:  "DB Update.Proc2",
			Start: testNow.Add(time.Minute * 5),
			End:   testNow.Add(time.Minute*5 + (time.Second * 15)),
		},
		{
			Name:  "Server Start",
			Start: testNow.Add(time.Minute * 4),
			End:   testNow.Add(time.Minute * 6),
		},
	},
	"2": []Event{
		{
			Name:  "a.b",
			Start: testNow,
			End:   testNow.Add(time.Minute),
		},
		{
			Name:  "a.c.d.e",
			Start: testNow.Add(time.Minute),
			End:   testNow.Add(time.Minute * 2),
		},
		{
			Name:  "a.c.d.f",
			Start: testNow.Add(time.Minute * 2),
			End:   testNow.Add(time.Minute * 3),
		},
	},
}

type TimelineResponse struct {
	*HalResponse
	First  time.Time `json:",omitempty"`
	Last   time.Time `json:",omitempty"`
	Events []Event   `json:",omitempty"`
}

func (tr *TimelineResponse) updateFirstLast() {
	for _, e := range tr.Events {
		if !tr.Last.After(e.End) {
			tr.Last = e.End
		}
	}

	tr.First = tr.Last

	for _, e := range tr.Events {
		if !tr.First.Before(e.Start) {
			tr.First = e.Start
		}
	}
}

func getTimeline(w rest.ResponseWriter, req *rest.Request) {
	rw := halResponseWriter{w}
	re := &TimelineResponse{HalResponse: MakeHalResponse()}

	context := req.PathParam("context")
	timeline := req.PathParam("timeline")

	if context == "test" {
		events, exists := testTimelines[timeline]

		if !exists {
			rest.NotFound(w, req)
			return
		} else {
			re.Events = events
			re.updateFirstLast()
			rw.WriteJSONResponse(req.URL, re)
		}
	} else {
		if !datastore.ContextExists(context) {
			rest.NotFound(w, req)
			return
		}

		events := datastore.GetTimelineEvents(context, timeline)
		for _, event := range events {
			re.Events = append(
				re.Events,
				Event{
					Name:  event.EventName,
					Host:  event.Host,
					Start: event.Start,
					End:   event.End,
				},
			)
		}

		re.updateFirstLast()
		rw.WriteJSONResponse(req.URL, re)
	}
}
