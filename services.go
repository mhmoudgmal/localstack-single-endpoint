package main

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// Services represents AWS services as localstack backends
type Services map[string]Backend

// DefaultLocalstackEndpoints returns the default localstack endpoints -
// mentioned in localstack: https://github.com/localstack/localstack#overview
func DefaultLocalstackEndpoints() (services Services) {
	content, readErr := ioutil.ReadFile("services.yml")
	if readErr != nil {
		log.Fatalf("[ERROR]: %v", readErr)
	}

	if parseErr := yaml.Unmarshal([]byte(content), &services); parseErr != nil {
		log.Fatalf("[ERROR]: %v", parseErr)
	}
	return
}

// Names returns the services names of the currently supported services.
func (services Services) Names() (servicesNames []string) {
	for serviceName := range services {
		servicesNames = append(servicesNames, serviceName)
	}
	return
}
