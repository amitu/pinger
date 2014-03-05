package pinger

import (
	"net"
	"fmt"
	"time"
	"net/http"
	"net/url"
	"encoding/json"
)

type Work struct {
	IP string
	HTTP, UDP string // only one of them need be "nonzero"
	Data string
}

type WorkResult struct {
	Success int
	Timeouts int
	Min, Max, Avg time.Duration	
	Data string
}

func (w Work) Execute() *WorkResult {
	result := w.FindPingStats()
	w.PostPingStats(result)
	return result
}

func (w Work) String() string {
	return fmt.Sprintf(
		"{http:%s udp:%s ip:%s data:%s}", w.HTTP, w.UDP, w.IP, w.Data,
	)
}

func (w Work) FindPingStats() *WorkResult {
	addr := net.ParseIP(w.IP)
	if addr == nil {
		fmt.Println("Invalid ip", w.IP)
		return nil
	}

	end := time.Now().Add(1e10)
	failures := 0
	success := 0
	var avgTtl, minTtl, maxTtl, sumTtls time.Duration
	minTtl = 1e12

	for {
		if time.Now().After(end) || (success + failures) >= 10 {
			break
		}

		ttl := Ping(&net.IPAddr{IP: addr}, 1e9)

		if ttl == nil {
			failures += 1
			continue
		}

		success += 1

		sumTtls += *ttl

		if *ttl > maxTtl {
			maxTtl = *ttl
		}

		if *ttl < minTtl {
			minTtl = *ttl
		}
	}

	if success > 0 {
		avgTtl = time.Duration(sumTtls / time.Duration(success))
	}

	result := WorkResult{
		Avg: avgTtl,
		Min: minTtl,
		Max: maxTtl,
		Success: success,
		Timeouts: failures,
	}

	return &result
}

func (w Work) PostPingStats(result *WorkResult) {
	if result == nil {
		fmt.Println("Result is nil")
		return
	}

	if w.HTTP != "" {
	    transport := http.Transport{
	        Dial: func(network, addr string) (net.Conn, error) {
	        	return net.DialTimeout(network, addr, 2e9)
	        },
	        ResponseHeaderTimeout: 2e9,
    	}
	    client := http.Client{
	        Transport: &transport,
	    }
		resp, err := client.PostForm(
			w.HTTP, url.Values{
				"Min": {result.Min.String()},
				"Max": {result.Max.String()}, 
				"Avg": {result.Avg.String()}, 
				"Success": {fmt.Sprintf("%d", result.Success)},
				"Timeouts": {fmt.Sprintf("%d", result.Timeouts)},
				"Data": {w.Data},
			},
		)
		if err != nil {
			fmt.Println(err)
			return
		}
		if resp.StatusCode >= 400 {
			fmt.Println("Worker: HTTP not ok:", resp.Status)
			return
		}
		fmt.Println("Worker: Posted data over HTTP.")
	} 

	if w.UDP != "" {
		conn, err := net.DialTimeout("udp", w.UDP, 1e9)
		if err != nil {
			fmt.Println("Error establishing connection to host: %s\n", err)
			return
		}

		result.Data = w.Data

		bytes, err := json.Marshal(result)
		if err != nil {
			fmt.Println("Error marshalling:", err)
			return
		}
		_, err = conn.Write(bytes)
		if err != nil {
			fmt.Println("Error marshalling:", err)
		}
		fmt.Println("Worker: Posted data over UDP.")
	}
}

