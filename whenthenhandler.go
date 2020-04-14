package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

func getHandleFunc(storage *InMemoryStore) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		then, err := storage.FindByRequest(NewStoreRequest(r.URL.Path, r.Method, Header(r.Header), r.Body))
		if errors.Is(err, NOT_FOUND) {
			w.WriteHeader(404)
			fmt.Fprintln(w, errors.Unwrap(err))
			return
		}
		if err != nil {
			w.WriteHeader(500)
			log.Println(err)

		}
		writeThen(w, then)
		return
	}
}

func writeThen(w http.ResponseWriter, then *Then) {
	if then.Delay > 0 {
		time.Sleep(time.Duration(then.Delay) * time.Millisecond)
	}
	for key, value := range then.Headers {
		w.Header().Set(key, "")
		for _, v := range value {
			w.Header().Add(key, v)
		}
	}
	w.WriteHeader(then.Status)
	fmt.Fprintln(w, then.Body)
}