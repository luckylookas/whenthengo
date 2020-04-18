package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

var NOT_FOUND = errors.New("")

type InMemoryStore map[string][]WhenThen

func (s InMemoryStore) getWhenThenKey(whenthen WhenThen) string {
	return fmt.Sprintf("%s#%s", CleanMethod(whenthen.When.Method), CleanUrl(whenthen.When.URL))
}

func (s InMemoryStore) getWhenThenKeyFromRequest(r StoreRequest) string {
	return fmt.Sprintf("%s#%s", CleanMethod(r.Method), CleanUrl(r.Url))
}

func (s InMemoryStore) Store(whenthen WhenThen) (key string, err error) {
	cleaned := WhenThen{
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

	if s[key] == nil {
		s[key] = []WhenThen{cleaned}
	} else {
		s[key] = append(s[key], cleaned)
	}
	return key, nil
}

func (s InMemoryStore) getByKey(key string) ([]WhenThen, error) {
	ret, ok := s[key]
	if !ok {
		return nil, NOT_FOUND
	}
	return ret, nil
}

func (s InMemoryStore) FindByRequest(storeRequest StoreRequest) (*Then, error) {
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

func (s InMemoryStore) headersMatch(storeRequest StoreRequest, candidate When) error {
	for key, value := range candidate.Headers {
		if !storeRequest.Headers.ContainsAllForKey(key, value...) {
			return fmt.Errorf("header mismatch for %s. %w", key, NOT_FOUND)
		}
	}
	return nil
}
