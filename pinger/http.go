package pinger

import (
	"net/http"
	"io"
	"fmt"
	"flag"
)

var (
	httpHostPort *string
)

func init() {
	httpHostPort = flag.String("http", "127.0.0.1:9918", "HTTP address.")
}

func pingitHandler(res http.ResponseWriter, req *http.Request) {
	job := Work {
		HTTP: req.FormValue("http"), 
		UDP: req.FormValue("udp"), 
		IP: req.FormValue("ip"),
		Data: req.FormValue("data"),
	}
	fmt.Printf("HTTP: %s.\n", job.String())
	workers <- job
	io.WriteString(res, "ok")
}

func statsHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Queue Size: %d", len(workers))
}

func SetupHTTPServer() {
	if (*httpHostPort == "disabled") {
		return
	}
	
	http.HandleFunc("/pingit", pingitHandler)
	http.HandleFunc("/stats", statsHandler)

	go http.ListenAndServe(*httpHostPort, nil)
	fmt.Printf("HTTP: Server started on %s.\n", *httpHostPort)
}