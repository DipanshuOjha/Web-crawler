package crawler

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

func Crawl(url string, depth int, visited *sync.Map, wg *sync.WaitGroup, linkchan chan<- string, sem chan struct{}) error {
	defer wg.Done()

	if depth <= 0 {
		return nil
	}

	if _, loaded := visited.LoadOrStore(url, true); loaded {
		return nil
	}

	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return fmt.Errorf("error while creatignn request %s : %v", url, err)
	}

	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("error while fetching detail %s : %s", url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status check out :- ,%s", resp.Status)
	}

	doc, err := html.Parse(resp.Body)

	if err != nil {
		return fmt.Errorf("failed to parse the html page")
	}
	var links []string
	var node func(n *html.Node, links *[]string)
	node = func(n *html.Node, links *[]string) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link := strings.TrimSpace(attr.Val)
					if link != "" && (strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://")) {
						*links = append(*links, link)
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			node(c, links)
		}
	}

	node(doc, &links)

	for _, link := range links {
		linkchan <- link
	}

	// setting limit till 5
	if len(links) > 5 {
		links = links[:5]
	}

	for _, link := range links {

		if _, loaded := visited.Load(link); !loaded {
			sem <- struct{}{}
			wg.Add(1)
			go func(link string) {
				Crawl(link, depth-1, visited, wg, linkchan, sem)
				<-sem
			}(link)
		}
	}

	return nil

}
