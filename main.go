package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
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
	// Output message
	fmt.Printf("\033[38;5;45mNew Request Received:\033[0m\n")

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

	// Otuput that we're listening
	fmt.Printf("\033[38;5;154mWebHook Debugger\033[0m by @mikemackintosh\n")
	fmt.Printf("Listening on port: %s\n", listenPort)
	fmt.Printf("%s\n\n", strings.Repeat("=", 24))

	// Start the server
	http.ListenAndServe(":"+listenPort, nil)
}
