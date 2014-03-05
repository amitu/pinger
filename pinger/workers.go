package pinger

import (
	"flag"
	"fmt"
	"sync"
)

var (
	nworkers *int
	wgWorkers sync.WaitGroup
	workers chan Work	
)

func init() {
	nworkers = flag.Int("workers", 10, "Number of workers.")
}

func SetupWorkers() {
	workers = make(chan Work, *nworkers * 2)

	for i := 0; i < *nworkers; i++ {
		go worker(i, workers)
	}
	fmt.Println("Workers: Started", *nworkers, "workers.")
}

func worker(id int, workers chan Work) {
	wgWorkers.Add(1)
	for {
		work := <- workers
		fmt.Printf("Worker %d: executing %s.\n", id, work.String())
		work.Execute()
	}
}