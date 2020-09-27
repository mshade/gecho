package main

import (
	"fmt"
	"net/http"
	"sort"
	s "strings"
)

func echo(w http.ResponseWriter, req *http.Request) {

	// Main attributes to return
	fmt.Fprintf(w, "%v: %v\n", "Host", req.Host)
	fmt.Fprintf(w, "%v: %v\n", "Remote Address", req.RemoteAddr)
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

	fmt.Fprintf(w, "%v", "\nHeaders\n")
	fmt.Fprintf(w, s.Join(allHeaders, ""))
}

func ip(w http.ResponseWriter, req *http.Request) {
	// Start with the plain client address.
	clientIP := s.Split(req.RemoteAddr, ":")[0]

	// Loop over header values we care about from least to most relevant.
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

	fmt.Fprintf(w, "%v\n", clientIP)
}

func main() {

	http.HandleFunc("/", echo)
	http.HandleFunc("/ip", ip)
	http.ListenAndServe(":8090", nil)

}
