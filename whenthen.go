package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func writeThen(w http.ResponseWriter, then *Then) {
	if then.Delay > 0 {
		time.Sleep(time.Duration(then.Delay) * time.Millisecond)
	}
	for key, value := range then.Headers {
		w.Header().Set(key, value)
	}
	w.WriteHeader(then.Status)
	fmt.Fprintln(w, then.Body)
}

func getHandleFunc(whenthens []*WhenThen) http.HandlerFunc {
	routes := map[string][]*WhenThen{}
	for _, whenthen := range whenthens {
		_, ok := routes[whenthen.When.URL]
		if !ok {
			routes[whenthen.When.URL] = []*WhenThen{}
		}
		routes[whenthen.When.URL] = append(routes[whenthen.When.URL], whenthen)
	}
	
	return func (w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		route, ok := routes[r.URL.Path]
		if !ok {
			log.Println("no whenthen for", r.URL.Path)
			w.WriteHeader(404)
			fmt.Fprintln(w, "no whenthen for", r.URL.Path)
			return
		}
		
		for _, whenthen := range route {

			if strings.ToLower(r.Method) != strings.ToLower(whenthen.When.Method) {
				log.Println("method mismatch")
				continue
			}
			matchedHeaders := true
			for key, value := range whenthen.When.Headers {
				stripped := stripWhenHeader(value)
				if !strings.Contains(stripWhenHeader(r.Header.Get(key)), stripped) {
					log.Println("header missmatch", key, value)
					matchedHeaders = false
					break;
				}
			}
			if !matchedHeaders {
				continue
			}
			
			requestBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println("error reading request", r.URL.Path)
				w.WriteHeader(500)
				fmt.Fprintln(w, "error reading request", r.URL.Path)
				return
			}
			
			if stripBody(string(requestBody)) != stripBody(whenthen.When.Body) {
				log.Println("body mismatch", stripBody(string(requestBody)), stripBody(whenthen.When.Body))
				continue
			}
			
			writeThen(w, &whenthen.Then)
			return
		}
		
		w.WriteHeader(404)
		fmt.Fprintln(w, "no whenthen matched method, path, headers and body for", r.URL.Path)
	}
}

func stripWhenHeader(value string) string {
	stripped := strings.ReplaceAll(value, ";", "")
	stripped = strings.ReplaceAll(stripped, ",", "")
	return strings.ToLower(stripped)
}

func stripBody(value string) string {
	stripped := strings.ReplaceAll(value, "\r\n", "")
	stripped = strings.ReplaceAll(stripped, "\n", "")
	stripped = strings.ReplaceAll(stripped, "\t", "")
	stripped = strings.ReplaceAll(stripped, " ", "")
	return stripped
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	configuration := LoadConfig()

	log.Println("configuration setup:", configuration)
	whenthens, err := Parse(configuration)
	if err != nil {
		log.Fatal(err)
		return
	}

	http.HandleFunc("/whenthengoup", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	http.HandleFunc("/", getHandleFunc(whenthens))
	
	go func () {
		log.Fatal(http.ListenAndServe(":"+strings.TrimPrefix(configuration.Port, ":"), nil))
	}()

	<-sigs
}
