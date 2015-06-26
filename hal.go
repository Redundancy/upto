package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"

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
	links Links,
	value interface{},
) error {
	rw.Header().Set("Content-Type", "application/hal+json")

	root := make(map[string]interface{})

	l := make(Links, len(links))
	for k, v := range links {
		l[k] = v
	}

	l["self"] = SimpleHREF("/api" + self.Path)

	// extract relevant fields from a struct
	typ := reflect.TypeOf(value)

	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	switch typ.Kind() {
	case reflect.Struct:
		val := reflect.ValueOf(value).Elem()
		for i := 0; i < typ.NumField(); i++ {
			fieldValue := val.Field(i)
			fieldType := typ.Field(i)

			if !fieldType.Anonymous && fieldValue.IsValid() {
				jsonTag := fieldType.Tag.Get("json")
				name := fieldType.Name

				switch {
				case jsonTag == "":
					// ignore
				case jsonTag == "-":
					// don't process it
					continue
				case strings.HasSuffix(jsonTag, omit):
					// actual omitting not implemented
					name = jsonTag[:len(jsonTag)-len(omit)]
				default:
					name = jsonTag
				}

				root[name] = fieldValue.Interface()
			}
		}
	default:
		return fmt.Errorf("Unimplemented type: %v", typ.Name())
	}

	root["_links"] = l

	return rw.WriteJson(root)
}
