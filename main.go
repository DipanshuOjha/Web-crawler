package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/DipanshuOjha/Web-crawler/crawler"
)

func main() {
	urlPtr := flag.String("url", "https://example.com", "Starting URL to crawl")
	depthPtr := flag.Int("depth", 2, "Crawl depth (non-negative integer)")
	concurrencyPtr := flag.Int("concurrency", 10, "Max concurrent goroutines")
	flag.Parse()

	if *depthPtr < 0 {
		fmt.Fprintln(os.Stderr, "Error: depth must be positive")
		os.Exit(1)
	}

	if *concurrencyPtr < 1 {
		fmt.Fprintln(os.Stderr, "Error: concurrency must be positive")
		os.Exit(1)
	}

	fmt.Printf("Starting to crawl %s (depth=%d, concurrency=%d)...\n", *urlPtr, *depthPtr, *concurrencyPtr)

	visited := &sync.Map{}
	uniquelink := &sync.Map{}
	var wg sync.WaitGroup
	linkchan := make(chan string, 100)

	links := []string{}
	sem := make(chan struct{}, *concurrencyPtr) // Semaphore: 10 concurrent goroutines to limit the gorotines
	start := time.Now()
	go func() {
		defer close(linkchan)
		wg.Add(1)
		go crawler.Crawl(*urlPtr, *depthPtr, visited, &wg, linkchan, sem)
		wg.Wait()
	}()

	for link := range linkchan {
		if _, loaded := uniquelink.LoadOrStore(link, true); !loaded {
			links = append(links, link)
		}
	}

	end := time.Now()
	diff := end.Sub(start).Seconds()

	fmt.Println("total time taken in secs before go rotines ", diff)

	if len(links) == 0 {
		fmt.Println("No links found on the page.")
		return
	}

	fmt.Printf("Found %d unique links on %s:\n", len(links), *urlPtr)
	for i, link := range links {
		fmt.Printf("%d. %s\n", i+1, link)
		time.Sleep(time.Millisecond * 300)
	}

}
