package main

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	proxyPort          string
	localstackHost     string
	defaultBackendPort string

	requestChannel     = make(chan *Request, 100)
	localstackServices = DefaultLocalstackEndpoints()
)

func main() {
	flag.StringVar(&proxyPort, "proxyPort", "9000",
		"the main application port")

	flag.StringVar(&defaultBackendPort, "defaultBackendPort", "9001",
		"the backend port for the application to fallback to in case no localstack backend found for a request",
	)

	flag.StringVar(&localstackHost, "localstackHost", "localhost",
		"the host where localstack is running and accessible")

	flag.Parse()

	go http.ListenAndServe(fmt.Sprintf(":%s", proxyPort),
		LocalstackSingleEndpoint{LocalstackHost: localstackHost})

	go http.ListenAndServe(fmt.Sprintf(":%s", defaultBackendPort), DefaultBackend{})

	for {
		select {
		case req := <-requestChannel:
			go func() {
				backend := BackendFor(req.Request, Backend{"", defaultBackendPort})
				forward(req, backend)
			}()
		}
	}
}

// LocalstackSingleEndpoint represents the the proxy server
type LocalstackSingleEndpoint struct {
	http.Handler
	LocalstackHost string
}

func (l LocalstackSingleEndpoint) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	done := make(chan bool, 1)
	request := &Request{
		LocalstackHost: l.LocalstackHost,
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
