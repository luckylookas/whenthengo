package main

import "github.com/JeremyLoy/config"

type Configuration struct {
	WhenThen string `config:"WHENTHEN"`
	Port     string `config:"PORT"`
}

func LoadConfig() *Configuration {
	conf := Configuration{
		WhenThen: "./whenthen.json",
		Port: "8080",
	}
	config.FromEnv().To(&conf)
	
	
	return &conf
}
