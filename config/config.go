package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	BrewfatherStreamURL *string
	BrewfatherUserId    *string
	BrewfatherApiKey    *string
}

func LoadConfig(cfg *Config) error {
	dat, err := ioutil.ReadFile("gobrewcfg.json")
	if err == nil {
		err = json.Unmarshal(dat, cfg)
	}

	if err == nil {
		debug, _ := json.Marshal(*cfg)
		log.Println("Loaded config :%s", string(debug))
	} else {
		log.Println("Failed to load config:")
		log.Println(err)
	}

	return err
}
