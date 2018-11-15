package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Request is the http request that will be forwarded to the downstream service.
type Request struct {
	*http.Request
	http.ResponseWriter

	Done chan bool // channel to notify when the request has finsihed.
}

func forward(req *Request, backend Backend) {
	reqURL, _ := url.Parse(req.Request.URL.String())

	reqURL.Scheme = "http"
	reqURL.Host = backend.String()

	newRequest, reqErr := http.NewRequest(
		req.Request.Method,
		reqURL.String(),
		req.Request.Body,
	)
	if reqErr != nil {
		fmt.Println("Failed creating new request", reqErr)
		req.ResponseWriter.WriteHeader(500)
		req.ResponseWriter.Write([]byte(http.StatusText(500)))
		req.Done <- true
		return
	}

	// Copy all headers from the origianl request to the new request
	for k, v := range req.Request.Header {
		headerVals := ""
		for _, headerVal := range v {
			headerVals += headerVal + " "
		}

		newRequest.Header.Add(k, headerVals)
	}

	fmt.Printf("Request at [ %s ] is forwarded to %s\n",
		req.Request.URL.String(), newRequest.Host)

	client := http.Client{}
	res, reqErr := client.Do(newRequest)

	if reqErr != nil {
		// FIXME: handle response properly, not all error are 500!
		req.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		req.Done <- true
		return
	}

	// Copy all the headers from the response to the original response writer.
	for k, v := range res.Header {
		headerVals := ""
		for _, headerVal := range v {
			headerVals += headerVal + " "
		}

		req.ResponseWriter.Header().Add(k, headerVals)
	}

	io.Copy(req.ResponseWriter, res.Body)

	req.Done <- true
}
