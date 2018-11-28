package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Request is the http request that will be forwarded to the downstream service.
type Request struct {
	LocalstackHost string

	*http.Request
	http.ResponseWriter

	Done chan bool // channel to notify when the request has finsihed.
}

func forward(req *Request, backend Backend) {
	reqURL, _ := url.Parse(req.Request.URL.String())

	reqURL.Scheme = "http"
	if req.LocalstackHost != "" {
		backend.Host = req.LocalstackHost
	}
	reqURL.Host = backend.String()

	reqBody, _ := ioutil.ReadAll(req.Request.Body)
	newRequest, newReqErr := http.NewRequest(req.Request.Method,
		reqURL.String(), bytes.NewBuffer(reqBody))

	if newReqErr != nil {
		log.Printf("[WARNING]: Failed forwarding the request to %s: %v",
			backend.String(), newReqErr)

		req.ResponseWriter.WriteHeader(http.StatusUnprocessableEntity)
		req.ResponseWriter.Write([]byte(newReqErr.Error()))
		req.Done <- true
		return
	}

	copyRequestHeaders(req.Request, newRequest)

	log.Printf("[INFO]: Request at [%s]{} is forwarded to [%s]",
		req.Request.URL.String(),
		newRequest.Host)

	client := http.Client{}
	res, reqErr := client.Do(newRequest)

	if reqErr != nil {
		req.ResponseWriter.WriteHeader(http.StatusBadGateway)
		req.ResponseWriter.Write([]byte(reqErr.Error()))
		req.Done <- true
		return
	}

	copyResponseHeaders(res, req)
	req.ResponseWriter.WriteHeader(res.StatusCode)

	io.Copy(req.ResponseWriter, res.Body)

	req.Done <- true
}

func copyRequestHeaders(originalRequest, newRequest *http.Request) {
	for k, v := range originalRequest.Header {
		newRequest.Header.Add(k, strings.Join(v, " "))
	}
}

func copyResponseHeaders(newResponse *http.Response, req *Request) {
	for k, v := range newResponse.Header {
		req.ResponseWriter.Header().Add(k, strings.Join(v, " "))
	}
}
