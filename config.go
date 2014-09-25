package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func LoadConfig(filename string) (NameValues, error) {
	result := make(NameValues)

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(content, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (nv NameValues) HasFields(fields []string) bool {
	result := true

	for _, field := range(fields) {
		_, found := nv[field]
		result = result && found
	}

	return result
}

const ConfigFile = "config.json"

var config NameValues
func init() {
	var err error
	config, err = LoadConfig(ConfigFile)
	if err != nil {
		log.Fatalf("Could not load config file %v because of error: %v!\n", ConfigFile, err)
	}

	requiredFields := []string{EndpointKey, UserKey, PasswordKey, SignatureKey}
	if !config.HasFields(requiredFields) {
		log.Fatalf("Required fields %v are missing from %v\n", requiredFields, ConfigFile)
	}
}
