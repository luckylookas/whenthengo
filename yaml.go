package main

import (
	"fmt"
	"github/luckylukas/whenthengo/types"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
)

type WhenThenYaml struct {
	When WhenYaml `yaml:"when"`
	Then ThenYaml `yaml:"then"`
}

type WhenYaml struct {
	Method  string                 `yaml:"method"`
	URL     string                 `yaml:"url"`
	Headers map[string]interface{} `yaml:"headers"`
	Body    string                 `yaml:"body"`
}

type ThenYaml struct {
	Status  int                    `yaml:"status"`
	Delay   int                    `yaml:"delay"`
	Headers map[string]interface{} `yaml:"headers"`
	Body    string                 `yaml:"body"`
}

type YamlParser struct {
}

func (_ YamlParser) String() string {
	return "yaml"
}

func (parser YamlParser) castHeaders(items map[string]interface{}) (headers map[string][]string) {
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

func (parser YamlParser) Parse(reader io.Reader) ([]*types.WhenThen, error) {
	whenthen := make([]*WhenThenYaml, 0)
	buffer, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(buffer, &whenthen)
	var ret = make([]*types.WhenThen, len(whenthen))
	for i, item := range whenthen {
		ret[i] = &types.WhenThen{
			When:types.When{
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
