package main

import (
	"errors"
	"fmt"
	"github/luckylukas/whenthengo/types"
	"log"
	"net/http"
	"time"
)

func getAddingFunc(storage *InMemoryStore) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		whenthens, err :=  JsonParser{}.Parse(r.Body)
		if err != nil {
			w.WriteHeader(500)
			log.Println(err)
			return
		}

		for _, item := range whenthens {
			_, err := storage.Store(*item)
			if err != nil {
				w.WriteHeader(500)
				log.Println(err)
				return
			}
		}
		w.WriteHeader(201)
	}
}

func getHandleFunc(storage *InMemoryStore) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		then, err := storage.FindByRequest(types.NewStoreRequest(r.URL.Path, r.Method, types.Header(r.Header), r.Body))
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

func writeThen(w http.ResponseWriter, then *types.Then) {
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