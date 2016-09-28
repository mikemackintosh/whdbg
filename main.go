package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
)

// listenPort is supplied to ListenAndServer
var listenPort string

func init() {
	// Port to listen via HTTP on
	//
	// #-> default is 8080
	flag.StringVar(&listenPort, "p", "8080", "Default listener port")
}

// handler will write the request dump to the response and stdout
func handler(w http.ResponseWriter, r *http.Request) {
	// generate the request dump
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		requestDump = []byte(err.Error())
	}

	// Output locally
	fmt.Printf(string(requestDump))

	// Output to response
	w.Write(requestDump)
}

func main() {
	// Parse flags
	flag.Parse()

	// Set the handler to output the request
	http.HandleFunc("/", handler)

	// Start the server
	http.ListenAndServe(":"+listenPort, nil)
}
