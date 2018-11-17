package main

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	// ProxyPort the main proxy port
	ProxyPort string
	// DefaultBackendPort the fallback backend port
	DefaultBackendPort string

	requestChannel     = make(chan *Request, 100)
	localstackServices = DefaultLocalstackEndpoints()
)

func main() {
	flag.StringVar(&ProxyPort, "ProxyPort", "9000",
		"the main application port")

	flag.StringVar(&DefaultBackendPort, "DefaultBackendPort", "9001",
		"the backend port for the application to fallback to in case no localstack backend found for a request",
	)

	flag.Parse()

	go http.ListenAndServe(fmt.Sprintf(":%s", ProxyPort), LocalstackSingleEndpoint{})
	go http.ListenAndServe(fmt.Sprintf(":%s", DefaultBackendPort), DefaultBackend{})

	for {
		select {
		case req := <-requestChannel:
			go func() {
				backend := BackendFor(req.Request, Backend{"", DefaultBackendPort})
				forward(req, backend)
			}()
		}
	}
}

// LocalstackSingleEndpoint represents the the proxy server
type LocalstackSingleEndpoint struct {
	http.Handler
}

func (LocalstackSingleEndpoint) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	done := make(chan bool, 1)
	request := &Request{
		ResponseWriter: res,
		Request:        req,
		Done:           done,
	}

	requestChannel <- request
	<-done
}

// DefaultBackend gets all requests that can't be handled by localstack
type DefaultBackend struct {
	http.Handler
}

func (DefaultBackend) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte("No localstack backend for this request"))
}
