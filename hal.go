package main

import (
	"encoding/json"
	"net/url"

	"github.com/ant0ine/go-json-rest/rest"
)

// SimpleHREF wraps a string href into a hal href object in json
type SimpleHREF string

// MarshalJSON serializes the uri as {"href": <uri>}
func (s SimpleHREF) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			HREF string `json:"href"`
		}{
			string(s),
		},
	)
}

type halResponseWriter struct {
	rest.ResponseWriter
}

// Links is a helper type to abbreviate map[string]interface{}
type Links map[string]interface{}

const omit = ",omitempty"

type HalResponse struct {
	Links map[string]interface{} `json:"_links,omitempty"`
}

func MakeHalResponse() *HalResponse {
	return &HalResponse{Links: make(map[string]interface{}, 1)}
}

func (h *HalResponse) SetSelf(url string) {
	if h.Links == nil {
		h.Links = make(map[string]interface{}, 1)
	}
	h.Links["self"] = SimpleHREF(url)
}

type HalResponseType interface {
	SetSelf(string)
}

/*
WriteJSONResponse provides a wrapper on JSON marshalling that
makes it easier to support some aspects of the HAL specification -
specifically embedding the _links field in the root of the object, alongside
other fields.

It has only the most basic support for json tags and may fail in all but the most
simple cases.
*/
func (rw halResponseWriter) WriteJSONResponse(
	self *url.URL,
	re HalResponseType,
) error {
	rw.Header().Set("Content-Type", "application/hal+json")
	re.SetSelf("/api" + self.Path)
	return rw.WriteJson(re)
}
