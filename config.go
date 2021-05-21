package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Configuration struct {
	Topic   string   `json:"topic"`   // kafka topic
	Group   string   `json:"groupID"` // kafka group
	Brokers []string `json:"brokers"` // kafka brokers
	Produce bool     `json:"produce"` // produce random messages
	SHA     string   `json:"sha"`     // sha version to use, e.g. sha1, sha256, sha512
}

// global variables
var Config Configuration

// String returns string representation of dbs Config
func (c *Configuration) String() string {
	data, err := json.Marshal(c)
	if err != nil {
		log.Println("ERROR: fail to marshal configuration", err)
		return ""
	}
	return string(data)
}

func ParseConfig(configFile string) error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println("unable to read config file", configFile, err)
		return err
	}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Println("unable to parse config file", configFile, err)
		return err
	}
	return nil
}
