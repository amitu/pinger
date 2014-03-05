package pinger 

import (
	"net"
	"fmt"
	"flag"
	"encoding/json"
)

var (
	udpHostPort *string
)

func init() {
	udpHostPort = flag.String("udp", "127.0.0.1:9918", "UDP address.")
}

func udpServer(conn net.PacketConn) {
	for {
		bytes := make([]byte, 512)
		n, _, err := conn.ReadFrom(bytes)

		if err != nil {
			fmt.Println(err)
			continue
		}

		bytes = bytes[:n]

		fmt.Println(bytes)

		var job Work
		err = json.Unmarshal(bytes, &job)

		if err != nil {
			fmt.Println(err)
		}

		if job.IP == "" {
			fmt.Println("UDP: got bad json, no IP supplied.", job.String())
			continue
		}
		
		workers <- job

		fmt.Printf("UDP: %s.\n", job.String())
	}
}

func SetupUDPServer() {
	if (*udpHostPort == "disabled") {
		return
	}

	conn, err := net.ListenPacket("udp4", *udpHostPort)
	if err != nil {
		panic(err)
	}

	go udpServer(conn)
	
	fmt.Printf("UDP: Server started on %s.\n", *udpHostPort)
}