package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/DipanshuOjha/Web-crawler/crawler"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <url>\n", os.Args[0])
		os.Exit(1)
	}

	time.Sleep(200 * time.Millisecond)
	fmt.Println("just got the link...")

	url := os.Args[1]

	depth := 2

	if len(os.Args) > 2 {
		var err error
		depth, err := strconv.Atoi(os.Args[2])
		if err != nil || depth < 0 {
			fmt.Println("Error: depth must be a non-negative integer")
			os.Exit(1)
		} else {
			time.Sleep(200 * time.Millisecond)
			fmt.Println("just got the depth......")
		}
	}

	fmt.Println("Starting to get your links.....")

	visited := &sync.Map{}
	var wg sync.WaitGroup
	linkchan := make(chan string, 100)

	links := []string{}
	start := time.Now()
	go func() {
		defer close(linkchan)
		wg.Add(1)
		go crawler.Crawl(url, depth, visited, &wg, linkchan)
		wg.Wait()
	}()

	for link := range linkchan {
		links = append(links, link)
	}

	end := time.Now()
	diff := end.Sub(start).Seconds()

	fmt.Println("total time taken in secs before go rotines ", diff)

	if len(links) == 0 {
		fmt.Println("No links found on the page.")
		return
	}

	fmt.Printf("Found %d links on %s:\n", len(links), url)
	for i, link := range links {
		fmt.Printf("%d. %s\n", i+1, link)
		time.Sleep(time.Millisecond * 300)
	}

}
