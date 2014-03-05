package main

// requires $ setcap cap_net_raw+ep /path/to/pinger or sudo

import (
	"github.com/amitu/pinger/pinger"
	"github.com/amitu/gutils"
	"flag"
)

func main() {
	flag.Parse()

	pinger.SetupWorkers()
	pinger.SetupICMP4Server()
	pinger.SetupHTTPServer()
	pinger.SetupUDPServer()

	gutils.WaitForCtrlC()
}