// Package jsonhal provides structs and methods to easily wrap your own data
// in a HAL compatible struct with support for hyperlinks and embedded resources
// HAL specification: http://stateless.co/hal_specification.html
package jsonhal

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

// Link represents a link in "_links" object
type Link struct {
	Href  string `json:"href"`
	Title string `json:"title,omitempty"`
}

// Embedded represents a resource in "_embedded" object
type Embedded interface{}

// Hal is used for composition, include it as anonymous field in your structs
type Hal struct {
	Links    map[string]*Link    `json:"_links,omitempty"`
	Embedded map[string]Embedded `json:"_embedded,omitempty"`
}

// SetLink sets a link (self, next, etc). Title argument is optional
func (h *Hal) SetLink(name, href, title string) {
	if h.Links == nil {
		h.Links = make(map[string]*Link, 0)
	}
	h.Links[name] = &Link{Href: href, Title: title}
}

// GetLink returns a link by name or error
func (h *Hal) GetLink(name string) (*Link, error) {
	if h.Links == nil {
		return nil, fmt.Errorf("Link \"%s\" not found", name)
	}
	link, ok := h.Links[name]
	if !ok {
		return nil, fmt.Errorf("Link \"%s\" not found", name)
	}
	return link, nil
}

// SetEmbedded adds a slice of objects under a named key in the embedded map
func (h *Hal) SetEmbedded(name string, embedded Embedded) {
	if h.Embedded == nil {
		h.Embedded = make(map[string]Embedded, 0)
	}
	h.Embedded[name] = embedded
}

// GetEmbedded returns a slice of embedded resources by name or error
func (h *Hal) GetEmbedded(name string) (Embedded, error) {
	if h.Embedded == nil {
		return nil, fmt.Errorf("Embedded \"%s\" not found", name)
	}
	embedded, ok := h.Embedded[name]
	if !ok {
		return nil, fmt.Errorf("Embedded \"%s\" not found", name)
	}
	return embedded, nil
}

// CountEmbedded counts number of embedded items
func (h *Hal) CountEmbedded(name string) (int, error) {
	e, err := h.GetEmbedded(name)
	if err != nil {
		return 0, err
	}
	if reflect.TypeOf(interface{}(e)).Kind() != reflect.Slice && reflect.TypeOf(interface{}(e)).Kind() != reflect.Map {
		return 0, errors.New("Embedded object is not a slice or a map")
	}
	return reflect.ValueOf(interface{}(e)).Len(), nil
}

// DecodeEmbedded decodes embedded objects into a struct
func (h *Hal) DecodeEmbedded(name string, result interface{}) error {
	e, err := h.GetEmbedded(name)
	if err != nil {
		return err
	}
	return mapstructure.Decode(interface{}(e), result)
}
