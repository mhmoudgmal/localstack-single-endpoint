package main

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// Services ....................
type Services map[string]Backend

func allServices() (services Services) {
	content, readErr := ioutil.ReadFile("services.yml")
	if readErr != nil {
		panic(readErr)
	}

	if parseErr := yaml.Unmarshal([]byte(content), &services); parseErr != nil {
		log.Fatalf("error: %v", parseErr)
	}
	return
}

func servicesNames() (servicesNames []string) {
	for service := range allServices() {
		servicesNames = append(servicesNames, service)
	}
	return
}
