package transport

import (
	"errors"
	"github.com/luckylukas/whenthengo/client"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func TestStorage_GetWhenThenKey(t *testing.T) {
	s := InMemoryStore{}
	assert.Equal(t, s.getWhenThenKey(client.WhenThen{
		When: client.When{
			Url:    "/path/123",
			Method: "GET",
		},
	}), "get#path/123")
}

func TestStorage_GetWhenThenKeyFromRequest(t *testing.T) {
	s := InMemoryStore{}
	assert.Equal(t, s.getWhenThenKeyFromRequest(StoreRequest{
		Url:    "/path/123",
		Method: http.MethodGet,
	}), "get#path/123")
}

func TestStorage_Store_and_Find_CleanableMismatches_Match(t *testing.T) {
	s := InMemoryStore{}
	test := client.WhenThen{
		When: client.When{
			Url:    "abc/def",
			Method: "get",
			Headers: CleanHeaders(map[string][]string{
				"accept": {"a", "b"},
			}),
			Body: "abc",
		},
		Then: client.Then{
			Status: 100,
		},
	}

	key, err := s.Store(test)
	assert.NoError(t, err)
	assert.NotNil(t, key)

	actual, err := s.FindByRequest(NewStoreRequest("/abc/def", "get", Header{
		"Accept": {"a", "b"},
	}, strings.NewReader("A b\nc")))

	assert.NoError(t, err)

	assert.Equal(t, test.Then, *actual)
}

func TestStorage_Store_and_Find_Header_KeyMismatch_NoMatch(t *testing.T) {
	s := InMemoryStore{}
	test := client.WhenThen{
		When: client.When{
			Url:    "abc/def",
			Method: "get",
		},
		Then: client.Then{
			Status: 100,
		},
	}

	key, err := s.Store(test)
	assert.NoError(t, err)
	assert.NotNil(t, key)

	_, err = s.FindByRequest(NewStoreRequest("/abc/xyz", "get", nil, nil))
	assert.Error(t, err)
	assert.True(t, errors.Is(err, NOT_FOUND))
}

func TestStorage_Store_and_Find_Header_ContentMismatch_NoMatch(t *testing.T) {
	s := InMemoryStore{}
	test := client.WhenThen{
		When: client.When{
			Url:     "abc/def",
			Method:  "get",
			Headers: map[string][]string{"something": {"this"}},
		},
		Then: client.Then{
			Status: 100,
			Delay:  1,
			Headers: map[string][]string{
				"accept": {"app/json", "app/xml"},
			},
			Body: "abc",
		},
	}

	key, err := s.Store(test)
	assert.NoError(t, err)
	assert.NotNil(t, key)

	_, err = s.FindByRequest(NewStoreRequest("/abc/def", "get", Header{"something": {"different"}}, nil))
	assert.Error(t, err)
	assert.True(t, errors.Is(err, NOT_FOUND))
}

func TestStorage_Store_and_Find_Header_When_IsSubsetOf_HeaderRequest_Match(t *testing.T) {
	s := InMemoryStore{}
	test := client.WhenThen{
		When: client.When{
			Url:     "abc/def",
			Method:  "get",
			Headers: map[string][]string{"something": {"this"}},
		},
		Then: client.Then{
			Status: 100,
			Delay:  1,
			Headers: map[string][]string{
				"accept": {"app/json", "app/xml"},
			},
			Body: "abc",
		},
	}

	key, err := s.Store(test)
	assert.NoError(t, err)
	assert.NotNil(t, key)

	v, err := s.FindByRequest(NewStoreRequest("/abc/def", "get", Header{"something": {"this", "different"}}, nil))
	assert.NoError(t, err)
	assert.NotNil(t, v)

}

func TestStorage_Store_and_Find_Header_Request_IsSubsetOf_HeaderWhen_NoMatch(t *testing.T) {
	s := InMemoryStore{}
	test := client.WhenThen{
		When: client.When{
			Url:     "abc/def",
			Method:  "get",
			Headers: map[string][]string{"something": {"this", "both"}},
		},
		Then: client.Then{
			Status: 100,
			Delay:  1,
			Headers: map[string][]string{
				"accept": {"app/json", "app/xml"},
			},
			Body: "abc",
		},
	}

	key, err := s.Store(test)
	assert.NoError(t, err)
	assert.NotNil(t, key)

	_, err = s.FindByRequest(NewStoreRequest("/abc/def", "get", Header{"something": {"this"}}, nil))
	assert.Error(t, err)
	assert.True(t, errors.Is(err, NOT_FOUND))
}

func TestStorage_Store_and_Find_Header_No_Intersection_Match(t *testing.T) {
	s := InMemoryStore{}
	test := client.WhenThen{
		When: client.When{
			Url:     "abc/def",
			Method:  "get",
			Headers: map[string][]string{"something": {"this", "both"}},
		},
		Then: client.Then{
			Status: 100,
			Delay:  1,
			Headers: map[string][]string{
				"accept": {"app/json", "app/xml"},
			},
			Body: "abc",
		},
	}

	key, err := s.Store(test)
	assert.NoError(t, err)
	assert.NotNil(t, key)

	v, err := s.FindByRequest(NewStoreRequest("/abc/def", "get", Header{"different": {"content"}}, nil))
	assert.NoError(t, err)
	assert.NotNil(t, v)
}

func TestStorage_Store_and_Find_Body_WhenBody_RequestNoBody_NoMatch(t *testing.T) {
	s := InMemoryStore{}
	test := client.WhenThen{
		When: client.When{
			Url:     "abc/def",
			Method:  "get",
			Headers: map[string][]string{"something": {"this", "both"}},
			Body:    "abc",
		},
		Then: client.Then{
			Status: 100,
		},
	}

	key, err := s.Store(test)
	assert.NoError(t, err)
	assert.NotNil(t, key)

	_, err = s.FindByRequest(NewStoreRequest("/abc/def", "get", nil, nil))
	assert.Error(t, err)
	assert.True(t, errors.Is(err, NOT_FOUND))
}

func TestStorage_Store_and_Find_Body_WhenNoBody_RequestBody_NoMatch(t *testing.T) {
	s := InMemoryStore{}
	test := client.WhenThen{
		When: client.When{
			Url:     "abc/def",
			Method:  "get",
			Headers: map[string][]string{"something": {"this", "both"}},
		},
		Then: client.Then{
			Status: 100,
		},
	}

	key, err := s.Store(test)
	assert.NoError(t, err)
	assert.NotNil(t, key)

	_, err = s.FindByRequest(NewStoreRequest("/abc/def", "get", nil, strings.NewReader("anc")))
	assert.Error(t, err)
	assert.True(t, errors.Is(err, NOT_FOUND))
}

func TestStorage_Store_and_Find_Body_Mismatch_NoMatch(t *testing.T) {
	s := InMemoryStore{}
	test := client.WhenThen{
		When: client.When{
			Url:     "abc/def",
			Method:  "get",
			Headers: map[string][]string{"something": {"this", "both"}},
			Body:    "abc",
		},
		Then: client.Then{
			Status: 100,
		},
	}

	key, err := s.Store(test)
	assert.NoError(t, err)
	assert.NotNil(t, key)

	_, err = s.FindByRequest(NewStoreRequest("/abc/def", "get", nil, strings.NewReader("anc")))
	assert.Error(t, err)
	assert.True(t, errors.Is(err, NOT_FOUND))
}

type FailReader struct {
}

var mock_error = errors.New("MOCK")

func (_ FailReader) Read([]byte) (n int, err error) {
	return 0, mock_error
}

func TestStorage_FindByRequest_BodyReader_Error(t *testing.T) {
	s := InMemoryStore{}
	test := client.WhenThen{
		When: client.When{
			Url:    "abc/def",
			Method: "get",
			Headers: CleanHeaders(map[string][]string{
				"accept": {"a", "b"},
			}),
			Body: "abc",
		},
		Then: client.Then{
			Status: 100,
		},
	}
	key, err := s.Store(test)
	assert.NoError(t, err)
	assert.NotNil(t, key)

	_, err = s.FindByRequest(NewStoreRequest("/abc/def", "get", Header{
		"Accept": {"a", "b"},
	}, FailReader{}))

	assert.Error(t, err)
	assert.True(t, errors.Is(err, mock_error))
}

// issue related

func TestStorage_Store_MultipleEntriesForKey_Issue11(t *testing.T) {
	s := InMemoryStore{}
	test1 := client.WhenThen{
		When: client.When{
			Url:    "abc/def",
			Method: "get",
			Headers: CleanHeaders(map[string][]string{
				"accept": {"a", "b"},
			}),
			Body: "abc",
		},
		Then: client.Then{
			Status: 100,
		},
	}
	key1, err := s.Store(test1)
	assert.NoError(t, err)
	assert.NotNil(t, key1)
	test2 := client.WhenThen{
		When: client.When{
			Url:    "abc/def",
			Method: "get",
			Headers: CleanHeaders(map[string][]string{
				"accept": {"a"},
			}),
			Body: "abc",
		},
		Then: client.Then{
			Status: 100,
		},
	}

	key2, err := s.Store(test2)
	assert.NoError(t, err)
	assert.NotNil(t, key2)
	assert.Equal(t, key1, key2)

	actual1, err1 := s.FindByRequest(NewStoreRequest("/abc/def", "get", Header{
		"accept": {"a"},
	}, strings.NewReader("A b\nc")))

	assert.NoError(t, err1)

	assert.Equal(t, test2.Then, *actual1)
}

func TestStorage_Store_MultipleConflictingEntriesForKey_Issue11_FindFirst(t *testing.T) {
	s := InMemoryStore{}
	test1 := client.WhenThen{
		When: client.When{
			Url:    "abc/def",
			Method: "get",
			Headers: CleanHeaders(map[string][]string{
				"accept": {"a", "b"},
			}),
			Body: "abc",
		},
		Then: client.Then{
			Status: 100,
		},
	}
	key1, err := s.Store(test1)
	assert.NoError(t, err)
	assert.NotNil(t, key1)
	test2 := client.WhenThen{
		When: client.When{
			Url:    "abc/def",
			Method: "get",
			Headers: CleanHeaders(map[string][]string{
				"accept": {"a"},
			}),
			Body: "abc",
		},
		Then: client.Then{
			Status: 100,
		},
	}

	key2, err := s.Store(test2)
	assert.NoError(t, err)
	assert.NotNil(t, key2)
	assert.Equal(t, key1, key2)

	//find the first one to match
	actual, err := s.FindByRequest(NewStoreRequest("/abc/def", "get", Header{
		"Accept": {"a"},
	}, strings.NewReader("A b\nc")))

	assert.NoError(t, err)

	assert.Equal(t, test1.Then, *actual)
}
