package main

import (
	"net/http"
)

var requestChannel = make(chan *Request, 100)

func main() {
	http.HandleFunc("/", reRequest)

	go http.ListenAndServe(":3000", nil)

	for {
		select {
		case req := <-requestChannel:
			backend := backendFor(req.Request)
			go forward(backend, req)
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
