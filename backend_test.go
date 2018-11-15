package main

import (
	"fmt"
	"net/http"
	"testing"

	"gotest.tools/assert"
)

func awsRequestLikeFor(serviceName string) (request *http.Request) {
	request, _ = http.NewRequest(
		"POST",
		"http://localhost:3000",
		nil,
	)

	credentialHeader := fmt.Sprintf(
		"Credential=AKIAIOSFODNN7EXAMPLE/20181114/us-east-1/%s/aws4_request",
		serviceName,
	)
	request.Header["Authorization"] = []string{credentialHeader}
	return
}

func TestBackend_String(t *testing.T) {
	backend := Backend{Host: "localhost", Port: "1111"}

	expected := "localhost:1111"
	got := backend.String()

	assert.Equal(t, expected, got)
}

func TestBackendFor_missingAuthorization(t *testing.T) {
	request, _ := http.NewRequest(
		"POST",
		"http://localhost:3000",
		nil,
	)

	expected := Backend{}
	got := BackendFor(request)

	assert.Equal(t, expected, got)
}

func TestBackendFor_inValidCredential(t *testing.T) {
	request, _ := http.NewRequest(
		"POST",
		"http://localhost:3000",
		nil,
	)

	request.Header["Authorization"] = []string{"Credential="}
	expected := Backend{}
	got := BackendFor(request)

	assert.Equal(t, expected, got)
}

func TestBackendFor_apigateway(t *testing.T) {
	request, _ := http.NewRequest(
		"POST",
		"http://localhost:3000/restapis/772571883A-Z/dev/_user_request_/blablabla",
		nil,
	)
	expected := Backend{Host: "localhost", Port: "4567"}
	got := BackendFor(request)

	assert.Equal(t, expected, got)
}

func TestBackendFor(t *testing.T) {
	services := Services{
		"s3":       Backend{Host: "localhost", Port: "4572"},
		"lambda":   Backend{Host: "localhost", Port: "4574"},
		"dynamodb": Backend{Host: "localhost", Port: "4569"},
		"kinesis":  Backend{Host: "localhost", Port: "4568"},
	}

	for serviceName, expectedBackend := range services {
		request := awsRequestLikeFor(serviceName)

		got := BackendFor(request)

		assert.Equal(t, expectedBackend, got)
	}
}
