package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

type MockSuccessStorage struct{}

func (_ MockSuccessStorage) Store(WhenThen) (string, error) {
	return "", nil
}

func (_ MockSuccessStorage) FindByRequest(_ StoreRequest) (*Then, error) {
	return &Then{
		Status: 201,
	}, nil
}


type MockFailStorage struct{
	Err error
}
func (m MockFailStorage) Store(WhenThen) (string, error) {
	if m.Err == nil {
		return "", errors.New("mock")
	}
	return "", m.Err
}

func (m MockFailStorage) FindByRequest(_ StoreRequest) (*Then, error) {
	if m.Err == nil {
		return nil, errors.New("mock")
	}
	return nil, m.Err
}

type MockSuccessParser struct{}

func (_ MockSuccessParser) Parse(reader io.Reader) ([]*WhenThen, error) {
	return []*WhenThen{{
		When: When{
			Method: "get",
			URL:    "/path",
		},
		Then: Then{
			Status: 200,
		},
	}}, nil

}

type MockFailParser struct{}

func (_ MockFailParser) Parse(reader io.Reader) ([]*WhenThen, error) {
	return nil, errors.New("some")
}

func TestParseAndStoreWhenThens_success(t *testing.T) {
	parsers["1"] = MockSuccessParser{}

	config := Configuration{
		WhenThen: fmt.Sprintf("%s%c%s.json", "test_resources", os.PathSeparator, t.Name()),
	}

	err := ParseAndStoreWhenThens(&config, &MockSuccessStorage{})
	assert.NoError(t, err)

	parsers["1"] = MockFailParser{}

	err = ParseAndStoreWhenThens(&config, &MockSuccessStorage{})
	assert.Error(t, err)
}


func TestParseAndStoreWhenThens_Invalid(t *testing.T) {
	parsers["1"] = MockSuccessParser{}

	config := Configuration{
		WhenThen: fmt.Sprintf("%s%c%s.json", "test_resources", os.PathSeparator, t.Name()),
	}

	err := ParseAndStoreWhenThens(&config, &MockSuccessStorage{})
	assert.NoError(t, err)

	parsers["1"] = MockFailParser{}

	err = ParseAndStoreWhenThens(&config, &MockSuccessStorage{})
	assert.Error(t, err)
}

func TestParseAndStoreWhenThens_noconfig(t *testing.T) {
	parsers["1"] = MockSuccessParser{}

	config := Configuration{
		WhenThen: "",
	}

	err := ParseAndStoreWhenThens(&config, &MockSuccessStorage{})
	assert.NoError(t, err)
}

func TestParseAndStoreWhenThens_noFile(t *testing.T) {
	parsers["1"] = MockSuccessParser{}

	config := Configuration{
		WhenThen: "/doesntexist",
	}

	err := ParseAndStoreWhenThens(&config, &MockSuccessStorage{})
	assert.NoError(t, err)
}

func TestValidate_valid(t *testing.T) {
	assert.NoError(t, Validate(&WhenThen{
		When: When{
			Method: "get",
			URL:    "/path",
		},
		Then: Then{
			Status: 200,
		},
	}))
}

func TestValidate_invalid(t *testing.T) {
	assert.Error(t, Validate(&WhenThen{
		When: When{
			Method: "",
			URL:    "/path",
		},
		Then: Then{
			Status: 200,
		},
	}))

	assert.Error(t, Validate(&WhenThen{
		When: When{
			Method: "get",
			URL:    "",
		},
		Then: Then{
			Status: 200,
		},
	}))

	assert.Error(t, Validate(&WhenThen{
		When: When{
			Method: "get",
			URL:    "/path",
		},
		Then: Then{
		},
	}))

	assert.Error(t, Validate(&WhenThen{

	}))
}
