package main

import (
	"encoding/json"
	"fmt"
	"github/luckylukas/whenthengo/types"
	"io"
	"io/ioutil"
)

type JsonParser struct {
}

type WhenThenJson struct {
	When WhenJson `json:"when"`
	Then ThenJson `json:"then"`
}

type WhenJson struct {
	Method  string                 `json:"method"`
	URL     string                 `json:"url"`
	Headers map[string]interface{} `json:"headers"`
	Body    string                 `json:"body"`
}

type ThenJson struct {
	Status  int                    `json:"status"`
	Delay   int                    `json:"delay"`
	Headers map[string]interface{} `json:"headers"`
	Body    string                 `json:"body"`
}

func (_ JsonParser) String() string {
	return "json"
}

func (parser JsonParser) castHeaders(items map[string]interface{}) (headers map[string][]string) {
	headers = make(map[string][]string)
	for k, v := range items {
		switch value := v.(type) {
		case string:
			headers[k] = []string{value}
		case []string:
			headers[k] = value
		case []interface{}:
			tmp := make([]string, len(value))
			for i, item := range value {
				tmp[i] = fmt.Sprintf("%v", item)
			}
			headers[k] = tmp
		default:
			headers[k] = []string{fmt.Sprintf("%v", value)}
		}
	}
	return headers
}

func (parser JsonParser) Parse(reader io.Reader) ([]*types.WhenThen, error) {
	whenthen := make([]*WhenThenJson, 0)
	buffer, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buffer, &whenthen)
	var ret = make([]*types.WhenThen, len(whenthen))
	for i, item := range whenthen {
		ret[i] = &types.WhenThen{
			When: types.When{
				Method:  item.When.Method,
				URL:     item.When.URL,
				Headers: parser.castHeaders(item.When.Headers),
				Body:    item.When.Body,
			},
			Then: types.Then{
				Status:  item.Then.Status,
				Delay:   item.Then.Delay,
				Headers: parser.castHeaders(item.Then.Headers),
				Body:    item.Then.Body,
			},
		}
	}
	return ret, err
}