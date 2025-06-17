package crawler

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func Crawl(url string, depth int, visited map[string]bool) ([]string, error) {

	if depth <= 0 {
		return []string{}, nil
	}

	if visited[url] {
		return []string{}, nil
	}

	visited[url] = true

	resp, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("error while fecthing the url ,%v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status check out :- ,%s", resp.Status)
	}

	doc, err := html.Parse(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to parse the html page")
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

	allLink := links

	for _, link := range links {

		if !visited[link] {
			childlinks, err := Crawl(link, depth-1, visited)

			if err != nil {
				fmt.Println("Error for exploring this link ", link, err)
				continue
			}

			allLink = append(allLink, childlinks...)
		}
	}

	return links, nil

}
