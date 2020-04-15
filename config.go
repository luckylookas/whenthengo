package main

import "github.com/JeremyLoy/config"

type Configuration struct {
	WhenThen string `config:"WHENTHEN"`
	Port     string `config:"PORT"`
}

func LoadConfig() *Configuration {
	conf := Configuration{
		WhenThen: "./whenthen.json",
		Port:     "80",
	}
	config.FromEnv().To(&conf)

	return &conf
}
