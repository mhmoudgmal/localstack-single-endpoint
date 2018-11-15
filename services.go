package main

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

func getSupportedServices() (services map[string]Backend) {
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
	for serviceName := range servicesBackends {
		servicesNames = append(servicesNames, serviceName)
	}
	return
}
