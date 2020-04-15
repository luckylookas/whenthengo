package main

import (
	"github.com/luckylukas/cleaningpipe"
	"io"
)

type WhenThen struct {
	When When
	Then Then
}

type When struct {
	Method  string
	URL     string
	Headers map[string][]string
	Body    string
}

type Then struct {
	Status  int
	Delay   int
	Headers map[string][]string
	Body    string
}

type Parser interface {
	Parse(reader io.Reader) ([]*WhenThen, error)
}

type Header map[string][]string

func (h Header) containsForKey(key string, value string) bool {
	if value == "" || h[key] == nil || len(h[key]) == 0 {
		return true
	}

	for _, i := range h[key] {
		if i == value {
			return true
		}
	}

	return false
}

func (h Header) ContainsAllForKey(key string, values ...string) (contains bool) {
	contains = true
	for _, v := range values {
		contains = contains && h.containsForKey(key, v)
	}
	return contains
}

type StoreRequest struct {
	Url     string
	Body    io.Reader
	Headers Header
	Method  string
}

func NewStoreRequest(url, method string, header Header, body io.Reader) StoreRequest {
	return StoreRequest{
		Url:     CleanUrl(url),
		Method:  CleanMethod(method),
		Body:    cleaningpipe.NewCleaningPipe(CleanBodyBytes, body),
		Headers: CleanHeaders(header),
	}
}

type Store interface {
	Store(WhenThen) (string, error)
	FindByRequest(StoreRequest) (*Then, error)
}
