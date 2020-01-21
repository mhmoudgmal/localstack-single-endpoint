package main

import (
	"sort"
	"testing"

	"gotest.tools/assert"
)

var services = Services{
	"s3":         Backend{Host: "localhost", Port: "4572"},
	"lambda":     Backend{Host: "localhost", Port: "4574"},
	"kinesis":    Backend{Host: "localhost", Port: "4568"},
	"dynamodb":   Backend{Host: "localhost", Port: "4569"},
	"apigateway": Backend{Host: "localhost", Port: "4567"},
	"sns":        Backend{Host: "localhost", Port: "4575"},
	"sqs":        Backend{Host: "localhost", Port: "4576"},
}

func TestDefaultLocalstackEndpointst(t *testing.T) {
	expectedServices := services
	got := DefaultLocalstackEndpoints()

	assert.DeepEqual(t, got, expectedServices)
}

func TestNames(t *testing.T) {
	expectedServicesNames := []string{"s3", "lambda", "kinesis", "dynamodb", "apigateway", "sns", "sqs"}
	got := services.Names()

	sort.Strings(got)
	sort.Strings(expectedServicesNames)

	assert.DeepEqual(t, got, expectedServicesNames)
}
