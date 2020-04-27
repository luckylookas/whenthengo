package transport

import (
	"errors"
	"fmt"
	"github.com/luckylukas/cleaningpipe"
	"github.com/luckylukas/whenthengo/client"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

var NOT_FOUND = errors.New("")

type InMemoryStore map[string][]client.WhenThen

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

func (s InMemoryStore) getWhenThenKey(whenthen client.WhenThen) string {
	return fmt.Sprintf("%s#%s", CleanMethod(whenthen.When.Method), CleanUrl(whenthen.When.Url))
}

func (s InMemoryStore) getWhenThenKeyFromRequest(r StoreRequest) string {
	return fmt.Sprintf("%s#%s", CleanMethod(r.Method), CleanUrl(r.Url))
}

func (s InMemoryStore) Store(whenthen client.WhenThen) (key string, err error) {
	cleaned := client.WhenThen{
		When: client.When{
			Method:  CleanMethod(whenthen.When.Method),
			Url:     CleanUrl(whenthen.When.Url),
			Headers: CleanHeaders(whenthen.When.Headers),
			Body:    CleanBodyString(whenthen.When.Body),
		},
		Then: client.Then{
			Status:  whenthen.Then.Status,
			Delay:   whenthen.Then.Delay,
			Headers: whenthen.Then.Headers,
			Body:    whenthen.Then.Body,
		},
	}

	key = s.getWhenThenKey(cleaned)

	if s[key] == nil {
		s[key] = []client.WhenThen{cleaned}
	} else {
		s[key] = append(s[key], cleaned)
	}
	return key, nil
}

func (s InMemoryStore) getByKey(key string) ([]client.WhenThen, error) {
	ret, ok := s[key]
	if !ok {
		return nil, NOT_FOUND
	}
	return ret, nil
}

func (s InMemoryStore) FindByRequest(storeRequest StoreRequest) (*client.Then, error) {
	key := s.getWhenThenKeyFromRequest(storeRequest)
	candidates, err := s.getByKey(key)
	if err != nil {
		return nil, err
	}

	buffer, err := ioutil.ReadAll(storeRequest.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading requestbody %s, %w", storeRequest.Url, err)
	}
	requestBody := string(buffer)

	for _, candidate := range candidates {
		if err = s.headersMatch(storeRequest, candidate.When); err != nil {
			log.Println(err.Error())
			continue
		}

		if strings.Compare(requestBody, candidate.When.Body) != 0 {
			log.Println("Body mismatch", requestBody, "!=", candidate.When.Body)
			continue
		}

		return &candidate.Then, nil
	}
	return nil, NOT_FOUND
}

func (s InMemoryStore) headersMatch(storeRequest StoreRequest, candidate client.When) error {
	for key, value := range candidate.Headers {
		if !storeRequest.Headers.ContainsAllForKey(key, value...) {
			return fmt.Errorf("header mismatch for %s. %w", key, NOT_FOUND)
		}
	}
	return nil
}
