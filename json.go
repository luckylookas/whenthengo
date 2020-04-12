package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type JsonParser struct {

}

func (_ JsonParser) String() string {
	return "json"
}

func (parser JsonParser) Parse(reader io.Reader) ([]*WhenThen, error) {
	whenthen := make([]*WhenThen, 0)
	buffer, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buffer, &whenthen)
	return whenthen, err
}
