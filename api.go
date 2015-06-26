package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
)

// TODO: put a shim REST API here to implement a
func init() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Get("/", rootHandler),
		rest.Get("/contexts", getContexts),
		rest.Get("/contexts/:context", getContext),
		rest.Get("/contexts/:context/:item", getContextItem),
	)

	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	mux.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
}

func rootHandler(w rest.ResponseWriter, req *rest.Request) {
	rw := halResponseWriter{w}
	rw.WriteJSONResponse(
		req.URL,
		Links{
			"contexts": SimpleHREF("./contexts"),
		},
		nil,
	)
}

func getContexts(w rest.ResponseWriter, req *rest.Request) {
	rw := halResponseWriter{w}
	rw.WriteJSONResponse(
		req.URL,
		Links{
			"items": []SimpleHREF{
				"./server_startup",
			},
		},
		nil,
	)
}

func getContext(w rest.ResponseWriter, req *rest.Request) {
	rw := halResponseWriter{w}
	rw.WriteJSONResponse(
		req.URL,
		Links{
			"latest": SimpleHREF("./1"),
			"items": []SimpleHREF{
				"./1",
			},
		},
		&struct {
			Name        string
			Description string
		}{
			"Server Startup Timing",
			"Events that occur during the startup of the super-duper server (v1)",
		},
	)
}

type Event struct {
	Name  string
	Host  string `json:",omitempty"`
	Start time.Time
	End   time.Time
}

func getContextItem(w rest.ResponseWriter, req *rest.Request) {
	rw := halResponseWriter{w}
	now := time.Now()
	rw.WriteJSONResponse(
		req.URL,
		Links{
			"context": SimpleHREF("./.."),
		},
		&struct {
			Events []Event
		}{
			Events: []Event{
				{
					Name:  "File Copy",
					Host:  "Megacity1",
					Start: now,
					End:   now.Add(time.Minute),
				},
				{
					Name:  "DB Update",
					Start: now,
					End:   now.Add(time.Minute * 5),
				},
				{
					Name:  "Server Start",
					Start: now.Add(time.Minute * 4),
					End:   now.Add(time.Minute * 6),
				},
			},
		},
	)
}
