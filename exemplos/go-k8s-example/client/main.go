package main

import (
	"flag"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	uri      = flag.String("u", "http://localhost:8080", "set url")
	clients  = flag.Int("c", 100, "number of clients")
	requests = flag.Int("r", 100, "set number of requests")
)

func main() {
	flag.Parse()

	wg := sync.WaitGroup{}

	for i := 0; i < *clients; i++ {
		wg.Add(1)
		go func(n int) {
			for j := 0; j < *requests; j++ {
				_, err := http.Get(*uri)
				if err != nil {
					log.Printf("Client [%02d] - Request Error: %s\n", n, err.Error())
					break
				} else {
					log.Printf("Client [%02d] - Request Success!\n", n)
				}

				time.Sleep(time.Second * 1)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	log.Println("Finished program")
}
