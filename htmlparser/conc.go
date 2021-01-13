package htmlparser

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Crawler is a type that stores crawler details
type Crawler struct {
	BaseURL string
	filter  chan Page
	queue   chan string
	visited map[string]bool
	inQueue map[string]bool
	mu      sync.RWMutex
}

var wg sync.WaitGroup

// NewCrawler will initialise a new crawler
// instance and starts to visit all pages
func NewCrawler(baseURL, url string) *Crawler {
	crwl := &Crawler{
		BaseURL: baseURL,
		queue:   make(chan string, 100),
		filter:  make(chan Page),
		visited: make(map[string]bool),
		inQueue: make(map[string]bool),
	}

	crwl.jobPool()

	crwl.queue <- url

	time.Sleep(2 * time.Second)
	wg.Wait()

	return crwl
}

// Crawl will crawl a given url
func (c *Crawler) Crawl(url string) {
	c.mu.Lock()
	c.visited[url] = true
	// c.inQueue[url] = false
	delete(c.inQueue, url)
	c.mu.Unlock()

	wg.Add(1)
	c.filter <- ParseURL(c.BaseURL, url)
	fmt.Println("Visited", url, len(c.visited), len(c.inQueue))
	wg.Done()
}

func (c *Crawler) jobPool() {
	go c.filterURLs()
	go c.sheduler()
}

func (c *Crawler) filterURLs() {
	for page := range c.filter {
		if page.Broken {
			fmt.Println("Broken", page.URL)
			wg.Done()
			continue
		}
		for _, link := range page.Links {
			url := cleanURL(link.URL)

			c.mu.RLock()
			if strings.Contains(url, c.BaseURL) && !c.visited[url] && !c.inQueue[url] {
				c.queue <- url
				c.inQueue[url] = true
			}
			c.mu.RUnlock()
		}
		wg.Done()
	}

}

func (c *Crawler) sheduler() {
	for url := range c.queue {
		wg.Add(1)
		go c.Crawl(url)
	}
}
