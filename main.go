package main

import (
	"crypto/hmac"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
)

// flagListPort is supplied to ListenAndServe.
var flagListPort string

// flagAuthtoken is used to perform authentication against whdbg.
var flagAuthtoken string

// flagFormat if supplied prints data in this format.
var flagFormat string

// flagResponseBody if supplied sets the response data for the request.
var flagResponseBody string

// flagStatusCode is used to set the response status code.
var flagStatusCode int

// Webhook is the json interprestaiton.
type Webhook map[string]interface{}

func init() {
	// Port to listen via HTTP on
	// #-> default is 8080
	flag.StringVar(&flagListPort, "p", "8080", "Default listener port")

	// Output Format
	// #-> default is empty
	flag.StringVar(&flagFormat, "f", "", "Request Output Format")

	// Auth Token for the Authorization: Bearer header
	// #-> default is empty
	flag.StringVar(&flagAuthtoken, "a", "", "Authentication header token")

	// Response body sets the expected response body for the request
	// #-> default is empty
	flag.StringVar(&flagResponseBody, "r", "", "Response body to use for the request")

	// StatusCode
	// #-> default is 200
	flag.IntVar(&flagStatusCode, "c", 200, "Response code set for requests")
}

// handler will write the request dump to the response and stdout
func handler(w http.ResponseWriter, r *http.Request) {
	if len(flagAuthtoken) > 0 {
		if _, ok := r.Header["Authorization"]; !ok {
			w.Header().Set("X-Whdbg-Err", "Missing authorization header")
			w.WriteHeader(flagStatusCode)
			w.Write([]byte("{\"status\":\"error\"}"))
			return
		}

		// Compare the two tokens
		if !hmac.Equal([]byte(r.Header.Get("Authorization")), []byte(flagAuthtoken)) {
			w.Header().Set("X-Whdbg-Err", "Incorrect authorization header")
			w.WriteHeader(flagStatusCode)
			w.Write([]byte("{\"status\":\"error\"}"))
			return
		}
	}

	// Output message
	fmt.Printf("\n\033[38;5;45mNew Request Received:\033[0m\n")

	// Check the content type header
	var requestDump []byte
	switch r.Header.Get("Content-type") {
	case "application/json":
		decodeJSON(w, r.Body)
	case "text/plain":
	case "text/csv":
	case "application/xml":
	case "text/xml":
	default:
		// generate the request dump
		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			requestDump = []byte(err.Error())
		}
		fmt.Printf(string(requestDump))
	}

	// if response body is set, print it first.
	if len(flagResponseBody) > 0 {
		w.WriteHeader(flagStatusCode)
		w.Write([]byte(flagResponseBody))
		return
	}

	// Output to response
	w.WriteHeader(flagStatusCode)
	w.Write(requestDump)
}

func decodeJSON(w http.ResponseWriter, body io.Reader) {
	var webhook Webhook

	decoder := json.NewDecoder(body)
	err := decoder.Decode(&webhook)
	if err != nil {
		fmt.Printf("\033[38;5;196mError:\033[0m %s", err)
		w.Header().Set("X-Whdbg-Err", "Missing authorization header")
		w.WriteHeader(flagStatusCode)
		w.Write([]byte("{\"status\":\"error\"}"))
		return
	}

	if len(flagFormat) > 0 {
		// Supply format string to formater
		fm, _ := NewFormater(flagFormat)
		// Parse webhook
		if err = fm.Parse(&webhook); err != nil {
			fmt.Printf("\033[38;5;196mError:\033[0m %s", err)
			w.Header().Set("X-Whdbg-Err", "Missing authorization header")
			w.WriteHeader(flagStatusCode)
			w.Write([]byte("{\"status\":\"error\"}"))
			return
		}
	} else {
		// Output locally
		b, _ := json.Marshal(webhook)
		fmt.Printf("%s", string(b))
	}
}

func main() {
	// Parse flags
	flag.Parse()

	// Set the handler to output the request
	http.HandleFunc("/", handler)

	// Otuput that we're listening
	fmt.Printf("\033[38;5;154mWebHook Debugger\033[0m by @mikemackintosh\n")
	fmt.Printf("Listening on port: %s\n", flagListPort)
	fmt.Printf("%s\n\n", strings.Repeat("=", 24))

	// Start the server
	http.ListenAndServe(":"+flagListPort, nil)
}
