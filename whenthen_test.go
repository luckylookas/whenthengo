package main

import (
	"fmt"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	port, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal(err)
	}

	config := Configuration{
		WhenThen: "",
		Port:     fmt.Sprintf("%d", port),
	}

	s := make(chan os.Signal)
	defer close(s)
	ready := make(chan struct{})
	defer close(ready)
	ret := make(chan error)
	defer close(ret)

	go func() {
		ret <- run(&config, InMemoryStore{}, s)
	}()
	<-time.After(1 * time.Second)
	s <- syscall.SIGINT

	assert.NoError(t, <-ret)

}
