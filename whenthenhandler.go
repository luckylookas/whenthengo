package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

func getAddingFunc(storage Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		whenthens, err := JsonParser{}.Parse(r.Body)
		if err != nil {
			w.WriteHeader(500)
			log.Println(err)
			return
		}

		for _, item := range whenthens {
			if Validate(item) != nil {
				w.WriteHeader(400)
				fmt.Fprint(w, "invalid json payload")
				return
			}
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

func getHandleFunc(storage Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			defer r.Body.Close()
		}
		then, err := storage.FindByRequest(NewStoreRequest(r.URL.Path, r.Method, Header(r.Header), r.Body))
		if errors.Is(err, NOT_FOUND) {
			w.WriteHeader(404)
			fmt.Fprint(w, errors.Unwrap(err))
			return
		}
		if err != nil {
			w.WriteHeader(500)
			log.Println(err)
			return
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
		w.Header().Del(key)
		for _, v := range value {
			w.Header().Add(key, v)
		}
	}
	w.WriteHeader(then.Status)
	fmt.Fprint(w, then.Body)
}
