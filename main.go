package main

import (
	"net/http"
)

var (
	requestChannel     = make(chan *Request, 100)
	localstackServices = DefaultLocalstackEndpoints()
)

func main() {
	http.HandleFunc("/", reRequest)

	go http.ListenAndServe(":3000", nil)

	for {
		select {
		case req := <-requestChannel:
			backend := BackendFor(req.Request)
			go forward(req, backend)
		}
	}
}

func reRequest(res http.ResponseWriter, req *http.Request) {
	done := make(chan bool)
	request := &Request{
		ResponseWriter: res,
		Request:        req,
		Done:           done,
	}

	requestChannel <- request
	<-done
}
