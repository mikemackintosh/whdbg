package main

import (
	"crypto/hmac"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
)

var (
	// flagListenPort is supplied to ListenAndServe.
	flagListenPort string

	// flagAuthtoken is used to perform authentication against whdbg.
	flagAuthtoken string

	// flagFormat if supplied prints data in this format.
	flagFormat string

	// flagResponseBody if supplied sets the response data for the request.
	flagResponseBody string

	// flagStatusCode is used to set the response status code.
	flagStatusCode int

	//go:embed web/build
	embeddedFiles embed.FS

	// Version of binary
	Version = "dev-20210803"

	// Create the hub
	hub = newHub()
)

// Webhook is the json interprestaiton.
type Webhook map[string]interface{}

func init() {
	// Port to listen via HTTP on
	// #-> default is 8080
	flag.StringVar(&flagListenPort, "p", "8080", "Default listener port")

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

// DebugRequest
type DebugPayload struct {
	Listener  string        `json:"listener"`
	Dump      string        `json:"dump"`
	URL       string        `json:"url"`
	UnixTime  int64         `json:"unixtimestamp"`
	Timestamp string        `json:"timestamp"`
	Request   *DebugRequest `json:"request"`
}

type DebugRequest struct {
	Host    string
	Method  string
	URL     *url.URL
	Header  http.Header
	Cookies []*http.Cookie

	Proto      string // "HTTP/1.0"
	ProtoMajor int    // 1
	ProtoMinor int    // 0

	ContentLength    int64
	PostForm         url.Values
	Form             url.Values
	MultipartForm    *multipart.Form
	TransferEncoding []string

	RemoteAddr string
	RequestURI string

	RequestURL string `json:"req_url"`

	Body string
}

func ReqToDbg(r *http.Request) *DebugRequest {
	var dbg DebugRequest = DebugRequest{
		Host:    r.Host,
		Method:  r.Method,
		URL:     r.URL,
		Header:  r.Header,
		Cookies: r.Cookies(),

		Proto:      r.Proto,
		ProtoMajor: r.ProtoMajor,
		ProtoMinor: r.ProtoMinor,

		ContentLength:    r.ContentLength,
		PostForm:         r.PostForm,
		Form:             r.Form,
		MultipartForm:    r.MultipartForm,
		TransferEncoding: r.TransferEncoding,

		RemoteAddr: r.RemoteAddr,
		RequestURI: r.RequestURI,
	}

	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		dbg.Body = err.Error()
	} else {
		dbg.Body = string(raw)
	}
	defer r.Body.Close()

	return &dbg
}

// handler will write the request dump to the response and stdout
func handler(w http.ResponseWriter, r *http.Request) {
	listener := strings.Replace(r.Host, ".whdbg.dev", "", -1)

	if _, ok := hub.subs[listener]; !ok {
		w.Write([]byte("Please create the listener first by visiting https://whdbg.dev/_/" + listener + "\n"))
		return
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

	// Gen the request output
	var requestDump []byte
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		requestDump = []byte(err.Error())
	}

	wh := DebugPayload{
		Listener:  listener,
		Request:   ReqToDbg(r),
		Timestamp: time.Now().Format(time.RFC1123),
		UnixTime:  time.Now().Unix(),
	}
	wh.URL = wh.Request.URL.String()

	wh.Dump = string(requestDump)

	b, err := json.Marshal(wh)
	if err != nil {
		requestDump = []byte(err.Error())
	}
	hub.subs[listener].hub.broadcast <- b

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

	// Output to response
	w.WriteHeader(hub.subs[listener].response.StatusCode)

	if !hub.subs[listener].response.Reflect {
		w.Write([]byte(hub.subs[listener].response.Body))
		return
	}

	w.Write(requestDump)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")
	if len(p) > 0 {
		if p[len(p)-1] == "update" {
			if err := updateResponse(p[len(p)-2], r); err != nil {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("{\"status\":\"error\"}"))
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{\"status\":\"ok\"}"))
			return
		}
	}
}

func updateResponse(listener string, r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	fmt.Println(string(body))
	if len(body) > 2 {
		var resp Response
		if err := json.Unmarshal(body, &resp); err != nil {
			log.Fatalf("-> UPD ERR: %s\n", err)
		}

		if len(resp.Body) == 0 {
			resp.Reflect = true
		} else {
			resp.Reflect = false
		}

		hub.subs[listener].response = resp
	}

	return nil
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

func getFileSystem() http.FileSystem {
	// Get the build subdirectory as the
	// root directory so that it can be passed
	// to the http.FileServer
	fsys, err := fs.Sub(embeddedFiles, "web/build")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

/*
request -> handler
		-> /ws/ --> WebSocket handler
		-> / --> serve page
		-> /_/fsdfsd -> serve page
*/

// hostDetectorMiddleware
func hostDetectorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Host, ".whdbg.dev") {
			http.FileServer(getFileSystem()).ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Parse flags
	flag.Parse()

	// Set the handler to output the request
	r := http.NewServeMux()

	// Create websocket listeners
	r.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		sub := r.URL.Path[len("/ws/"):]
		fmt.Printf("=> Starting WS listener: %s\n", sub)
		serveWs(hub, sub, w, r)
	})

	r.HandleFunc("/api/", apiHandler)

	r.HandleFunc("/_/", func(w http.ResponseWriter, r *http.Request) {
		d, _ := embeddedFiles.ReadFile("web/build/index.html")
		w.WriteHeader(http.StatusOK)
		w.Write(d)
	})

	r.Handle("/", hostDetectorMiddleware(http.HandlerFunc(handler)))

	// Otuput that we're listening
	fmt.Printf("\033[38;5;154mWebHook Debugger\033[0m by @mikemackintosh\n")
	fmt.Printf("%s\n\n", strings.Repeat("=", 24))

	var wg sync.WaitGroup
	fmt.Printf("Starting Websocket Server...\n")
	go func(wg *sync.WaitGroup) {
		hub.run()
		defer wg.Done()
	}(&wg)
	wg.Add(1)

	fmt.Printf("Listening HTTP on port %s...\n", flagListenPort)
	go func(wg *sync.WaitGroup) {
		log.Fatal(http.ListenAndServe(":"+flagListenPort, r))
		defer wg.Done()
	}(&wg)
	wg.Add(1)

	wg.Wait()
}
