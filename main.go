package main

import (
	"crypto/hmac"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

var (
	// flagListPort is supplied to ListenAndServe.
	flagListPort string

	// flagAuthtoken is used to perform authentication against whdbg.
	flagAuthtoken string

	// flagFormat if supplied prints data in this format.
	flagFormat string

	// flagResponseBody if supplied sets the response data for the request.
	flagResponseBody string

	// flagStatusCode is used to set the response status code.
	flagStatusCode int

	// Create the hub
	hub = newHub()
)

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
	fmt.Printf("r.URL.Host = %s\n", r.URL.Host)
	fmt.Printf("r.Header.Get(Host) = %s\n", r.Header.Get("Host"))
	fmt.Printf("r.Host = %s\n", r.Host)
	fmt.Printf("r.URL.Hostname() = %s\n", r.URL.Hostname())

	if r.Host == "whdbg.dev" {
		if r.URL.Path == "/" {

			home.Execute(w, map[string]interface{}{
				"subs": hub.subs,
			})
			return
		}

		if strings.Contains(r.URL.Path, "/_/") {
			sub := r.URL.Path[len("/_/"):]
			if len(sub) > 0 {
				sock.Execute(w, map[string]string{"sub": sub, "ws": "wss://" + r.Host + "/ws/" + sub})
				return
			}
		}
	}

	if len(flagAuthtoken) > 0 {
		if _, ok := r.Header["Authorization"]; !ok {
			w.Header().Set("X-Whdbg-Err", "Missing authorization header")
			w.WriteHeader(flagStatusCode)
			w.Write([]byte("{\"status\":\"error\"}"))
			return
		}

		// Compare the two tokens
		if !hmac.Equal([]byte(r.Header.Get("Authorization")), []byte("Bearer "+flagAuthtoken)) {
			w.Header().Set("X-Whdbg-Err", "Incorrect authorization header")
			w.WriteHeader(flagStatusCode)
			w.Write([]byte("{\"status\":\"error\"}"))
			return
		}
	}

	// Output message
	fmt.Printf("\n\033[38;5;45mNew Request Received:\033[0m\n")

	// Gen the request output
	var requestDump []byte
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		requestDump = []byte(err.Error())
	}

	listener := strings.Replace(r.Host, ".whdbg.dev", "", -1)
	if ws, ok := hub.subs[listener]; ok {
		ws.hub.broadcast <- requestDump
	} else {
		w.Write([]byte("Please create the listener first by visiting https://whdbg.dev/_/" + listener + "\n"))
		return
	}

	// listener := strings.Replace(r.Host, ".whdbg.dev", "", -1)
	// push(w, r, listener, requestDump)

	// Check the content type header
	switch r.Header.Get("Content-type") {
	case "application/json":
		decodeJSON(w, r.Body)
		//	case "text/plain":
		//	case "text/csv":
		//	case "application/xml":
		//	case "text/xml":
	default:
		// generate the request dump
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

	go hub.run()

	// Set the handler to output the request
	http.HandleFunc("/", handler)

	// Create websocket listeners
	http.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		sub := r.URL.Path[len("/ws/"):]
		serveWs(hub, sub, w, r)
	})

	// Otuput that we're listening
	fmt.Printf("\033[38;5;154mWebHook Debugger\033[0m by @mikemackintosh\n")
	fmt.Printf("Listening on port: %s\n", flagListPort)
	fmt.Printf("%s\n\n", strings.Repeat("=", 24))

	// Start the server
	if len(os.Getenv("CERT")) > 0 && len(os.Getenv("KEY")) > 0 {
		if httpErr := http.ListenAndServeTLS(":"+flagListPort, os.Getenv("CERT"), os.Getenv("KEY"), nil); httpErr != nil {
			log.Fatal("The process exited with https error: ", httpErr.Error())
		}
	}

	http.ListenAndServe(":"+flagListPort, nil)
}
