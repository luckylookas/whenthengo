package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type WhenThen struct {
	When When
	Then Then
}

type When struct {
	Method  string
	URL     string
	Headers map[string][]string
	Body    string
}

type Then struct {
	Status  int
	Delay   int
	Headers map[string][]string
	Body    string
}

type WhenThenClient struct {
	uri     string
	backlog []WhenThen
	lock    sync.RWMutex
}

type WBuilder struct {
	when   When
	client *WhenThenClient
}

type TBuilder struct {
	when   When
	then   Then
	client *WhenThenClient
}

func (b *WBuilder) WithMethod(method string) *WBuilder {
	b.when.Method = method
	return b
}

func (b *WBuilder) WithUri(uri string) *WBuilder {
	b.when.URL = uri
	return b
}

func (b *WBuilder) WithHeader(key, value string) *WBuilder {
	b.when.Headers[key] = append(b.when.Headers[key], value)
	return b
}

func (b *WBuilder) ClearHeadersForKey(key string) *WBuilder {
	b.when.Headers[key] = []string{}
	return b
}

func (b *WBuilder) ClearHeaders() *WBuilder {
	b.when.Headers = make(map[string][]string)
	return b
}

func (b *WBuilder) WithBody(body string) *WBuilder {
	b.when.Body = body
	return b
}

func (b *WBuilder) ThenReply() *TBuilder {
	t := TBuilder{
		client: b.client,
		when:   b.when,
		then: Then{
			Status:  200,
			Headers: make(map[string][]string),
		},
	}
	return &t
}

func (t *TBuilder) WithBody(body string) *TBuilder {
	t.then.Body = body
	return t
}

func (t *TBuilder) WithStatus(status int) *TBuilder {
	t.then.Status = status
	return t
}

func (t *TBuilder) WithDelay(delay int) *TBuilder {
	t.then.Delay = delay
	return t
}

func (t *TBuilder) WithHeader(key, value string) *TBuilder {
	t.then.Headers[key] = append(t.then.Headers[key], value)
	return t
}

func (t *TBuilder) ClearHeadersForKey(key string) *TBuilder {
	t.then.Headers[key] = []string{}
	return t
}

func (t *TBuilder) ClearHeaders() *TBuilder {
	t.then.Headers = map[string][]string{}
	return t
}

func (t *TBuilder) AndDo() TerminatingClient {
	t.client.add(WhenThen{
		When: t.when,
		Then: t.then,
	})
	return t.client
}

func (t *TBuilder) And() *WhenThenClient {
	t.client.add(WhenThen{
		When: t.when,
		Then: t.then,
	})
	return t.client
}

func (c *WhenThenClient) add(wt WhenThen) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.backlog = append(c.backlog, wt)
}

func NewClient(host, port string) (*WhenThenClient) {
	return &WhenThenClient{
		uri: fmt.Sprintf("http://%s:%s/whenthengo/whenthen", host, strings.TrimLeft(port, ":")),
	}
}

func (c *WhenThenClient) WhenRequest() *WBuilder {
	return &WBuilder{
		client: c,
		when: When{
			Method:  "get",
			URL:     "/",
			Headers: make(map[string][]string),
		},
	}
}

func (c *WhenThenClient) Publish(ctx context.Context) error {
	c.lock.Lock()
	body, err := json.Marshal(c.backlog)
	if err != nil {
		c.lock.Unlock()
		return err
	}
	c.backlog = make([]WhenThen, 0)
	c.lock.Unlock()

	req, err := http.NewRequest(http.MethodPost, c.uri, bytes.NewReader(body))
	if err != nil {
		return err
	}

	done := make(chan *http.Response)
	ers := make(chan error)

	defer close(done)
	defer close(ers)

	go func() {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			ers <- err
		}
		done <- resp
	}()

	select {
	case resp := <-done:
		if resp.StatusCode != 201 {
			return errors.New("error publishing to whenthengo, status " + resp.Status)
		}
		return nil
	case err = <-ers:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *WhenThenClient) Return() []WhenThen {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.backlog
}

type TerminatingClient interface {
	Publish(ctx context.Context) error
	Return() []WhenThen
}
