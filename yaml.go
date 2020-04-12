package main

import (
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
)

type YamlParser struct {

}

func (_ YamlParser) String() string {
	return "yaml"
}

func (parser YamlParser) Parse(reader io.Reader) ([]*WhenThen, error) {
	whenthen := make([]*WhenThen, 0)
	buffer, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(buffer, &whenthen)
	return whenthen, nil
}
