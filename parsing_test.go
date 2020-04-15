package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/luckylukas/whenthengo/types"
	"io"
	"os"
	"testing"
)

type Mockstorage struct{}

func (_ Mockstorage) Store(types.WhenThen) (string, error) {
	return "", nil
}

func (_ Mockstorage) FindByRequest(_ types.StoreRequest) (*types.Then, error) {
	return nil, nil
}

type MockSuccessParser struct{}

func (_ MockSuccessParser) Parse(reader io.Reader) ([]*types.WhenThen, error) {
	return []*types.WhenThen{}, nil
}

type MockFailParser struct{}

func (_ MockFailParser) Parse(reader io.Reader) ([]*types.WhenThen, error) {
	return nil, errors.New("some")
}

func TestParseAndStoreWhenThens_success(t *testing.T) {
	parsers["1"] = MockSuccessParser{}

	config := Configuration{
		WhenThen: fmt.Sprintf("%s%c%s.json", "test_resources", os.PathSeparator, t.Name()),
	}

	err := ParseAndStoreWhenThens(&config, &Mockstorage{})
	assert.NoError(t, err)

	parsers["1"] = MockFailParser{}

	err = ParseAndStoreWhenThens(&config, &Mockstorage{})
	assert.Error(t, err)
}
