package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	var storage InMemoryStore = make(map[string]*WhenThen)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	configuration := LoadConfig()

	if err := ParseAndStoreWhenThens(configuration, storage); err != nil {
		log.Fatal(err)
		return
	}

	http.HandleFunc("/whenthengo/up", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	http.HandleFunc("/", getHandleFunc(&storage))

	go func() {
		log.Fatal(http.ListenAndServe(":"+strings.TrimPrefix(configuration.Port, ":"), nil))
	}()

	<-sigs
}
