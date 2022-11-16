package pkg

import (
	"github.com/BurntSushi/toml"
	"log"
)

type Config struct {
	Port     string `toml:"port"`
	Host     string `toml:"db_host"`
	DBPort   string `toml:"db_port"`
	User     string `toml:"db_user"`
	Password string `toml:"db_password"`
	Name     string `toml:"db_name"`
}

func NewConfig() *Config {
	config := &Config{}
	_, err := toml.DecodeFile("config.toml", config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}
