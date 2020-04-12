package main

import (
	"errors"
	"io"
	"log"
	"os"
)

type WhenThen struct {
	When When
	Then Then
}

type When struct {
	Method  string                 `yaml:"method",json:"method"`
	URL     string                 `yaml:"url",json:"url"`
	Headers map[string]string `yaml:"headers",json:"headers"`
	Body    string                 `yaml:"body",json:"body"`
}

type Then struct {
	Status  int                    `yaml:"status"json:"status"`
	Delay   int                    `yaml:"delay"json:"delay"`
	Headers map[string]string `yaml:"headers"json:"headers"`
	Body    string                 `yaml:"body"json:"body"`
}

type Parser interface {
	Parse(reader io.Reader) ([]*WhenThen, error)
}

var parsers = map[string]Parser{}

func init() {
	parsers["json"] = JsonParser{}
	parsers["yaml"] = YamlParser{}
}

func Parse (configuration *Configuration) (ret []*WhenThen, err error) {
	var file *os.File
	for _, parser := range parsers {
		file, err = os.Open(configuration.WhenThen)
		if err != nil {
			return nil, err
		}
		ret, err = parser.Parse(file)
		if err != nil {
			log.Println("parsing with ", parser, "failed")
			log.Println(err)
		}
		if err == nil && ret != nil {
			return ret, nil
		}
	}
	return nil, errors.New("SORRY! no parser could parse the contents of whenthen file.")
}