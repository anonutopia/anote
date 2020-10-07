package main

import (
	"encoding/json"
	"log"
	"os"
)

// Config struct holds all our configuration
type Config struct {
	BotName         string   `json:"bot_name"`
	WavesNodeAPIKey string   `json:"waves_node_api_key"`
	NodeAddress     string   `json:"node_address"`
	Dev             bool     `json:"dev"`
	Debug           bool     `json:"debug"`
	TelegramAPIKey  string   `json:"telegram_api_key"`
	InitialPrice    uint64   `json:"initial_price"`
	EmailAddress    string   `json:"email_address"`
	TokenID         string   `json:"token_id"`
	FounderAddress  string   `json:"founder_address"`
	Hostname        string   `json:"hostname"`
	ShoutTime       int      `json:"shout_time"`
	PostgreSQL      string   `json:"postgre_sql"`
	NodeHost        string   `json:"node_host"`
	ShoutAddress    string   `json:"shout_address"`
	Exclude         []string `json:"exclude"`
	AintID          string   `json:"aint_id"`
}

// Load method loads configuration file to Config struct
func (sc *Config) Load(configFile string) error {
	file, err := os.Open(configFile)

	if err != nil {
		log.Printf("[Config.Load] Got error while opening config file: %v", err)
		return err
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&sc)

	if err != nil {
		log.Printf("[Config.Load] Error while decoding JSON: %v", err)
		return err
	}

	return nil
}

func initConfig() *Config {
	c := &Config{}
	c.Load("config.json")
	return c
}
