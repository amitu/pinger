package pinger 

import (
	"time"
	"net"
	"fmt"
	"math/rand"
	"sync"
)

// so there is a ping "server" server runs in a goroutine, and can
// have only one instance server gets request to send ping over one
// channel, and sends ping packets over another.

// the server can be communicated using apis defined here, server
// channels are available as global variables in this package

// there is only one api, ping(ip, deadline), which blocks, and
// returns either nil if no ping received in the given deadline, or
// time.Duration, the rtt (round trip time).

// server always sends 64bytes of data.

var (
	pingChan chan PingJob
	jobMap map[int] PingJob
	jobMapMutex sync.Mutex
)

type PingJob struct {
	Reply chan *time.Duration
	Deadline time.Time
	Addr *net.IPAddr
}

func getUniqueJobID() int {
	for {
		id := rand.Intn(0xffff)
		if _, ok := jobMap[id]; !ok {
			return id
		}
	}
}

func (j PingJob)ping(conn *net.IPConn) {
	jobMapMutex.Lock()
	defer jobMapMutex.Unlock()

	id := getUniqueJobID()
	seq := rand.Intn(0xffff)
	jobMap[id] = j
	sendPing(conn, j.Addr, id, seq)
}

func SetupICMP4Server() {
	rand.Seed(time.Now().UnixNano())
	pingChan = make(chan PingJob)
	jobMap = make(map[int] PingJob)

	go icmpServer()
	go idCleaner()	
}

func idCleaner() {
	for {
		time.Sleep(1e9)
		now := time.Now()

		jobMapMutex.Lock()
		for id, job := range jobMap {
			if job.Deadline.Before(now) {
				delete(jobMap, id)
			}
		}
		jobMapMutex.Unlock()
	}
}

func icmpReader(conn *net.IPConn) {
	for {
		reply, err := readOnePingReply(conn, time.Now().Add(1e9))

		if err != nil {
			continue
		} 

		// fmt.Println("icmpReader: got reply:", reply)

		jobMapMutex.Lock()

		if job, ok := jobMap[reply.id]; ok {
			job.Reply <- &reply.rtt
			delete(jobMap, reply.id)
		} else {
			fmt.Println("Droppping packet with id =", reply.id)
		}

		jobMapMutex.Unlock()
	}
}

func icmpServer() {
	conn, err := createICMP4Connection()
	if err != nil {
		fmt.Println(err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	go icmpReader(conn)

	for {
		job := <- pingChan
		job.ping(conn)
	}
}

func Ping(addr *net.IPAddr, duration time.Duration) *time.Duration {
	reply := make(chan *time.Duration, 3)
	pingChan <- PingJob{
		Reply: reply, 
		Addr: addr, 
		Deadline: time.Now().Add(duration),
	}
	select {
	case rtt := <- reply:
		if rtt == nil {
			fmt.Printf("Ping: %s: got timeout.\n", addr.String())
		} else {
			fmt.Printf("Ping: %s: %s.\n", addr.String(), rtt)
		}
		return rtt
	case <- time.After(duration):
		fmt.Printf("Ping: %s: no reply.\n", addr.String())
		return nil
	}
}