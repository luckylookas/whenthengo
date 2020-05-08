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
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	configuration := LoadConfig()
	err := run(configuration, InMemoryStore{}, sigs)
	if err != nil {
		log.Fatal(err)
	}
}

func run(configuration *Configuration, storage Store, signals <-chan os.Signal) error {

	if err := ParseAndStoreWhenThens(configuration, storage); err != nil {
		return err
	}

	http.HandleFunc("/whenthengo/whenthen", getAddingFunc(storage))
	http.HandleFunc("/whenthengo/up", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	http.HandleFunc("/", getHandleFunc(storage))

	go func() {
		log.Fatal(http.ListenAndServe(":"+strings.TrimPrefix(configuration.Port, ":"), nil))
	}()

	<-signals
	return nil
}
