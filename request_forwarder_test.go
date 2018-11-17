package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"gotest.tools/assert"
)

func credentialFor(serviceName string) string {
	return fmt.Sprintf("Credential=AKIAIOSFODNN7EXAMPLE/20181117/us-east-1/%s/aws4_request", serviceName)
}

// A simplified version of the <LocalstackSingleEndpoint> http handler
func reRequest(res http.ResponseWriter, req *http.Request) {
	done := make(chan bool)
	reReq := &Request{
		ResponseWriter: res,
		Request:        req,
		Done:           done,
	}

	backend := BackendFor(reReq.Request, Backend{"", "9001"})
	go forward(reReq, backend)
	<-done
}

func TestForward_noLocalstackBackendDetected(t *testing.T) {
	go http.ListenAndServe(":9001", DefaultBackend{})
	time.Sleep(150 * time.Millisecond)

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
