package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// for the config.json
type Config struct {
	CompanyListPath string   `json:"company_list_csv_path"`
	JobListPath     string   `json:"job_list_csv_path"`
	Keywords        []string `json:"keywords"`
}

/*
readConfigFile

Reads the config.json file and will return it as a
config struct.

Will not return an error object. Will abort the program if the
file cannot be processed correctly.
*/
func ReadConfigFile() Config {
	file, err := os.Open("./config.json")
	if err != nil {
		log.Fatal("error: Could not load config.json file")
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("error: Had trouble reading config.json as bytes")
	}

	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatal("error: Could not unmarshall config.json")
	}

	return config
}
