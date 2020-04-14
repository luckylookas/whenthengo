package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

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
		Url:     cleanUrl(url),
		Method:  cleanMethod(method),
		Body:    CleanBodyPipe{body},
		Headers: cleanHeaders(header),
	}
}

type Store interface {
	Store(WhenThen) (string, error)
	FindByRequest(StoreRequest) (*Then, error)
}

type InMemoryStore map[string]*WhenThen

var NOT_FOUND = errors.New("")

func (s InMemoryStore) getWhenThenKey(whenthen *WhenThen) string {
	return fmt.Sprintf("%s#%s", cleanMethod(whenthen.When.Method), cleanUrl(whenthen.When.URL))
}

func (s InMemoryStore) getWhenThenKeyFromRequest(r StoreRequest) string {
	return fmt.Sprintf("%s#%s", cleanMethod(r.Method), cleanUrl(r.Url))
}

func (s InMemoryStore) Store(whenthen WhenThen) (key string, err error) {
	cleaned := &WhenThen{
		When{
			Method:  cleanMethod(whenthen.When.Method),
			URL:     cleanUrl(whenthen.When.URL),
			Headers: cleanHeaders(whenthen.When.Headers),
			Body:    cleanBodyString(whenthen.When.Body),
		},
		Then{
			Status:  whenthen.Then.Status,
			Delay:   whenthen.Then.Delay,
			Headers: whenthen.Then.Headers,
			Body:    whenthen.Then.Body,
		},
	}

	key = s.getWhenThenKey(cleaned)
	s[key] = cleaned
	return key, nil
}

func (s InMemoryStore) getByKey(key string) (*WhenThen, error) {
	ret, ok := s[key]
	if ! ok {
		return nil, NOT_FOUND
	}
	return ret, nil
}

func (s InMemoryStore) FindByRequest(r StoreRequest) (*Then, error) {
	key := s.getWhenThenKeyFromRequest(r)
	item, err := s.getByKey(key)
	if err != nil {
		return nil, err
	}

	for key, value := range item.When.Headers {
		if !r.Headers.ContainsAllForKey(key, value...) {
			return nil, fmt.Errorf("no whenthen for header values %s=%s, %w", key, value, NOT_FOUND)
		}
	}

	requestBody, err := ioutil.ReadAll(CleanBodyPipe{r.Body})
	if err != nil {
		return nil, fmt.Errorf("error reading request %s, %w", r.Url, err)
	}

	if strings.Compare(string(requestBody), item.When.Body) != 0 {
		log.Println("Body mismatch", string(requestBody), item.When.Body)
		return nil, fmt.Errorf("no whenthen for request Body %w", NOT_FOUND)
	}

	return &item.Then, nil
}
