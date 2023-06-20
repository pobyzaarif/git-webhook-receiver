package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type AppSetting struct {
	AppHost              string `json:"app_host"`
	AppPort              string `json:"app_port"`
	AppBasicAuthUsername string `json:"app_basic_auth_username"`
	AppBasicAuthPassword string `json:"app_basic_auth_password"`
}

type Mapping struct {
	RepoName   string `json:"repo_name"`
	BranchName string `json:"branch_name"`
	Command    string `json:"command"`
}

type WebhookSetting struct {
	Github struct {
		Mapping []Mapping `json:"mapping"`
	} `json:"github"`
	Gitlab struct {
		Mapping []Mapping `json:"mapping"`
	} `json:"gitlab"`
	Bitbucket struct {
		Mapping []Mapping `json:"mapping"`
	} `json:"bitbucket"`
}

type Config struct {
	AppSetting     AppSetting     `json:"app_setting"`
	WebhookSetting WebhookSetting `json:"webhook_setting"`
}

func LoadConfig(filename string) *Config {
	// Read the JSON file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to read config file: %v", err))
	}

	// Unmarshal the JSON data into the config struct
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to unmarshal config JSON: %v", err))
	}

	return &config
}
