package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

var NOT_FOUND = errors.New("")

type InMemoryStore map[string]*WhenThen

func (s InMemoryStore) getWhenThenKey(whenthen *WhenThen) string {
	return fmt.Sprintf("%s#%s", CleanMethod(whenthen.When.Method), CleanUrl(whenthen.When.URL))
}

func (s InMemoryStore) getWhenThenKeyFromRequest(r StoreRequest) string {
	return fmt.Sprintf("%s#%s", CleanMethod(r.Method), CleanUrl(r.Url))
}

func (s InMemoryStore) Store(whenthen WhenThen) (key string, err error) {
	cleaned := &WhenThen{
		When: When{
			Method:  CleanMethod(whenthen.When.Method),
			URL:     CleanUrl(whenthen.When.URL),
			Headers: CleanHeaders(whenthen.When.Headers),
			Body:    CleanBodyString(whenthen.When.Body),
		},
		Then: Then{
			Status:  whenthen.Then.Status,
			Delay:   whenthen.Then.Delay,
			Headers: whenthen.Then.Headers,
			Body:    whenthen.Then.Body,
		},
	}

	key = s.getWhenThenKey(cleaned)
	log.Println("adding when for ", key)
	s[key] = cleaned
	return key, nil
}

func (s InMemoryStore) getByKey(key string) (*WhenThen, error) {
	ret, ok := s[key]
	if !ok {
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

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading request %s, %w", r.Url, err)
	}

	if strings.Compare(string(requestBody), item.When.Body) != 0 {
		log.Println("Body mismatch", string(requestBody), item.When.Body)
		return nil, fmt.Errorf("no whenthen for request Body %w", NOT_FOUND)
	}

	return &item.Then, nil
}
