package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	types "github.com/luckylukas/whenthengo/types"
	"net/http"
	"strings"
	"sync"
)

type WhenThenClient struct {
	uri     string
	backlog []types.WhenThen
	lock    sync.RWMutex
}

type WBuilder struct {
	when   types.When
	client *WhenThenClient
}

type TBuilder struct {
	when   types.When
	then   types.Then
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
	b.when.Headers = make(types.Header)
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
		then: types.Then{
			Status:  200,
			Headers: make(types.Header),
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
	t.then.Headers = types.Header{}
	return t
}

func (t *TBuilder) AndDo() TerminatingClient {
	t.client.add(types.WhenThen{
		When: t.when,
		Then: t.then,
	})
	return t.client
}

func (t *TBuilder) And() *WhenThenClient {
	t.client.add(types.WhenThen{
		When: t.when,
		Then: t.then,
	})
	return t.client
}

func (c *WhenThenClient) add(wt types.WhenThen) {
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
		when: types.When{
			Method:  "get",
			URL:     "/",
			Headers: make(types.Header),
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
	c.backlog = make([]types.WhenThen, 0)
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

func (c *WhenThenClient) Return() []types.WhenThen {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.backlog
}

type TerminatingClient interface {
	Publish(ctx context.Context) error
	Return() []types.WhenThen
}
