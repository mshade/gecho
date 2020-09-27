package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	s "strings"

	"github.com/gorilla/handlers"
)

// echo any interesting attributes about the request back to the client
func echo(w http.ResponseWriter, req *http.Request) {
	// main attributes to return
	fmt.Fprintf(w, "%v: %v\n", "Host", req.Host)
	fmt.Fprintf(w, "%v: %v\n", "Remote Address", req.RemoteAddr)
	fmt.Fprintf(w, "%v: %v\n", "Request Method", req.Method)
	fmt.Fprintf(w, "%v: %v\n", "Request URI", req.RequestURI)

	// Assemble a slice of headers
	var allHeaders []string
	for name, headers := range req.Header {
		for _, entry := range headers {
			allHeaders = append(allHeaders, fmt.Sprintf("%v: %v\n", name, entry))
		}
	}

	// sort our headers for consistency and readability
	sort.Strings(allHeaders)

	// spit it all out.
	fmt.Fprintf(w, "%v", "\nHeaders\n")
	fmt.Fprintf(w, s.Join(allHeaders, ""))
}

// return the client IP, as best as we can determine.
func ip(w http.ResponseWriter, req *http.Request) {
	// start with the plain client address.
	clientIP := s.Split(req.RemoteAddr, ":")[0]

	// loop over header values we care about from least to most relevant.
	// the sort order of these is undefined.
	for name, headers := range req.Header {
		if s.Contains(name, "X-Forwarded-For") {
			clientIP = s.Split(headers[0], ",")[0]
		} else if s.Contains(name, "X-Real-Ip") {
			clientIP = headers[0]
		} else if s.Contains(name, "Cf-Connecting-Ip") {
			clientIP = headers[0]
			break
		}
	}

	// gotcha!
	fmt.Fprintf(w, "%v\n", clientIP)
}

func health(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "ok\n")
}

func main() {

	// allow port override from env
	port := "8090"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	fmt.Println("listening on port " + port)

	http.HandleFunc("/", echo)
	http.HandleFunc("/ip", ip)

	// Healthcheck - don't log reqs
	http.Handle("/_healthz", handlers.CombinedLoggingHandler(ioutil.Discard, http.HandlerFunc(health)))

	http.ListenAndServe(":"+port, handlers.CombinedLoggingHandler(os.Stdout, http.DefaultServeMux))

}
