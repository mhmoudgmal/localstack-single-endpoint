package main

import (
	"fmt"
	"gotest.tools/assert"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func credentialFor(serviceName string) string {
	return fmt.Sprintf("Credential=AKIAIOSFODNN7EXAMPLE/20181117/us-east-1/%s/aws4_request", serviceName)
}

func setupDefaultBackend() {
	go http.ListenAndServe(fmt.Sprintf(":%s", defaultBackendPort), DefaultBackend{})
}

// A simplified version of the <LocalstackSingleEndpoint> http handler
func reRequest(res http.ResponseWriter, req *http.Request) {
	done := make(chan bool)
	reReq := &Request{
		ResponseWriter: res,
		Request:        req,
		Done:           done,
	}

	backend := BackendFor(reReq.Request)
	go forward(reReq, backend)
	<-done
}

func TestForward_noLocalstackBackendDetected(t *testing.T) {
	setupDefaultBackend()

	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(reRequest)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("wrong status code returned: got %v want %v",
			status, http.StatusBadRequest)
	}

	expectedBody := "No localstack backend for this request"
	assert.Equal(t, expectedBody, rr.Body.String())
}

func TestForward_supportedLocalstackBackends(t *testing.T) {
	services := []string{"s3", "lambda", "apigateway", "dynamodb", "kinesis"}

	expectedBodyRegx := `^Post http://.*:\d{4}/: dial tcp .*:\d{4}: .*: connection refused$`
	r := regexp.MustCompile(expectedBodyRegx)

	for _, service := range services {
		req, err := http.NewRequest("POST", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("Authorization", credentialFor(service))

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(reRequest)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadGateway {
			t.Errorf("wrong status code returned: got %v want %v",
				status, http.StatusBadRequest)
		}

		if !r.MatchString(rr.Body.String()) {
			t.Errorf("response body differes: expected body to match (%v) got(%v)",
				expectedBodyRegx, rr.Body.String())
		}
	}
}
