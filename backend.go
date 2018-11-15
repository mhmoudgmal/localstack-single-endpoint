package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var serviceNamesRegx = fmt.Sprintf(`%s`,
	strings.Join(localstackServices.Names(), "|"))

const (
	regionRegx       = `(?:eu|us)\-(?:east|west|central)\-\d{1}`
	dateRegx         = `\d{4}\d{2}\d{2}`
	accessKeyIDRegex = `[A-Z0-9]{20}`

	// Localstack apigateway regex
	apigatewayURLRegx = `/restapis/[A-Za-z0-9_\-]+/[A-Za-z0-9_\-]+/_user_request_/.*`
)

// Backend represents the location where requests are forwarded to.
type Backend struct {
	Host string `yaml:"Host"`
	Port string `yaml:"Port"`
}

func (backend Backend) String() string {
	return backend.Host + ":" + backend.Port
}

// BackendFor analyzes an http request and returns the corresponding
// localstack endpoint according to AWS docs:
// docs.aws.amazon.com/AmazonS3/latest/API/sigv4-auth-using-authorization-header.html)
//
// [Authorization] header should contain info about the service.
// according to Credential component:
//	  <your-access-key-id>/<date>/<aws-region>/<aws-service>/aws4_request
//
// except for the REST apis that should be forwarded to apigateway endpoint
func BackendFor(req *http.Request) (backend Backend) {
	apigatewayRgx := regexp.MustCompile(apigatewayURLRegx)

	if apigatewayRgx.MatchString(req.URL.String()) {
		backend = localstackServices["apigateway"]
		return
	}

	var authorizationHeader []string
	var found bool

	if authorizationHeader, found = req.Header["Authorization"]; !found {
		log.Println("[WARNING]: Authorization header is missing")
		return
	}

	credentialFmt := fmt.Sprintf(`Credential=(%s)/(%s)/(%s)/(%s)/aws4_request`,
		accessKeyIDRegex,
		dateRegx,
		regionRegx,
		serviceNamesRegx,
	)

	credentialHeaderRgx := regexp.MustCompile(credentialFmt)

	matchedCredentialRgx := credentialHeaderRgx.FindStringSubmatch(
		authorizationHeader[0],
	)

	if len(matchedCredentialRgx) < 5 {
		log.Println("[WARNING]: Credential header does not match AWS's format")
		return
	}

	backend = localstackServices[matchedCredentialRgx[4]]
	return
}
