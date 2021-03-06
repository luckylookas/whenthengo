package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

var parsers = map[string]Parser{}

func init() {
	parsers["json"] = JsonParser{}
	parsers["yaml"] = YamlParser{}
}

func Validate(whenthen *WhenThen) error {
	if whenthen.When.Method == "" || whenthen.When.URL == "" || whenthen.Then.Status == 0 {
		return errors.New("a whenthen requires at least a url, a method and a response status to work")
	}
	return nil
}

func ParseAndStoreWhenThens(configuration *Configuration, storage Store) error {
	log.Println("loading config from", configuration.WhenThen)
	if configuration.WhenThen == "" {
		log.Println("no configuration, starting empty")
		return nil
	}
	for _, parser := range parsers {
		file, err := os.Open(configuration.WhenThen)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}

		items, err := parser.Parse(file)
		if err != nil {
			log.Println("parsing with ", parser, "failed")
			log.Println(err)
			continue
		}

		for _, item := range items {
			if err := Validate(item); err != nil {
				return err
			}
			if key, err := storage.Store(*item); err != nil {
				return fmt.Errorf("could not store whenthen for %s: %v", key, err)
			}
		}
		return nil
	}
	return errors.New("SORRY! no parser could parse contents of whenthen file.")
}
