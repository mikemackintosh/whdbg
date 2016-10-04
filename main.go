package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
)

// listenPort is supplied to ListenAndServer
var listenPort string

// format if supplied prints data in this format
var format string

// Webhook is the json interprestaiton
type Webhook map[string]interface{}

func init() {
	// Port to listen via HTTP on
	//
	// #-> default is 8080
	flag.StringVar(&listenPort, "p", "8080", "Default listener port")

	// Output Format
	//
	// #-> default is empty
	flag.StringVar(&format, "f", "", "Request Output Format")
}

// handler will write the request dump to the response and stdout
func handler(w http.ResponseWriter, r *http.Request) {
	var webhook Webhook

	// Output message
	fmt.Printf("\n\033[38;5;45mNew Request Received:\033[0m\n")

	// generate the request dump
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		requestDump = []byte(err.Error())
	}

	// If we supply a webhook format, output the data in this format
	if len(format) > 0 {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&webhook)
		if err != nil {
			fmt.Printf("> Error: %s", err)
		}

		// Supply format string to formater
		fm, _ := NewFormater(format)
		// Parse webhook
		fm.Parse(&webhook)
	} else {
		// Output locally
		fmt.Printf(string(requestDump))
	}

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
