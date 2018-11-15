package main

import (
	"fmt"
	"net/http"
)

const defaultBackendPort = "9000"

var (
	requestChannel     = make(chan *Request, 100)
	localstackServices = DefaultLocalstackEndpoints()
)

func main() {
	go http.ListenAndServe(":3000", LocalstackSingleEndpoint{})
	go http.ListenAndServe(fmt.Sprintf(":%s", defaultBackendPort), DefaultBackend{})

	for {
		select {
		case req := <-requestChannel:
			go func() {
				backend := BackendFor(req.Request)
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
