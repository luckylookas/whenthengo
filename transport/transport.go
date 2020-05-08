package transport

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type WhenThenGoTrasport struct {
	store InMemoryStore
}

func newTransport(config []WhenThen) WhenThenGoTrasport {
	wt := WhenThenGoTrasport{
		store: InMemoryStore{},
	}
	for _, i := range config {
		wt.store.Store(i)
	}
	return wt
}

func (wt WhenThenGoTrasport) RoundTrip(request *http.Request) (*http.Response, error) {
	then, err := wt.store.FindByRequest(NewStoreRequest(request.URL.Path, request.Method, Header(request.Header), request.Body))
	return writeThen(then), err
}

func writeThen(then *Then) *http.Response {
	w := &http.Response{}
	if then == nil {
		w.StatusCode = 404
		w.Status = http.StatusText(404)
		return w
	}
	if then.Delay > 0 {
		time.Sleep(time.Duration(then.Delay) * time.Millisecond)
	}
	for key, value := range then.Headers {
		w.Header.Del(key)
		for _, v := range value {
			w.Header.Add(key, v)
		}
	}
	w.StatusCode = then.Status
	w.Status = http.StatusText(then.Status)
	w.Body = ioutil.NopCloser(strings.NewReader(then.Body))
	return w
}

func NewWhenThenGoHttpClient(config []WhenThen) *http.Client {
	return &http.Client{
		Transport: newTransport(config),
	}

}
