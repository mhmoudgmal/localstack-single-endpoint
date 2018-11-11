package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// Backend represents the location where requests are forwarded to.
type Backend struct {
	Host string `yaml:"Host"`
	Port string `yaml:"Port"`
}

func (backend Backend) String() string {
	return backend.Host + ":" + backend.Port
}

func backendFor(req *http.Request) (backend Backend) {
	// restapi? -> apigateway
	if strings.Contains(req.URL.String(), "restapi") {
		backend = allServices()["apigateway"]
		return
	}

	var authorizationHeader []string
	var ok bool

	// (aws-sdk || awscli)? -> Authorization header should contain info about the service.
	// according to Credential component:
	//	  <your-access-key-id>/<date>/<aws-region>/<aws-service>/aws4_request
	//
	// docs.aws.amazon.com/AmazonS3/latest/API/sigv4-auth-using-authorization-header.html
	if authorizationHeader, ok = req.Header["Authorization"]; !ok {
		fmt.Println("Authorization header is missing.")
		return
	}

	supportedServices := strings.Join(servicesNames(), "|")
	credentialHeader := fmt.Sprintf(`Credential=.*/*/(%s)/`, supportedServices)

	rgx, err := regexp.Compile(credentialHeader)
	if err != nil {
		fmt.Println("Credential header does not match AWS's format.")
		return
	}

	service := rgx.FindStringSubmatch(authorizationHeader[0])[1]
	backend = allServices()[service]
	return
}
