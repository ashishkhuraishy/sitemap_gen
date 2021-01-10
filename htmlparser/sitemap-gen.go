package htmlparser

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// SiteMap is a struct which stores all
// the info about the website
type SiteMap struct {
	Domain      string
	PageCount   int
	BrokenLinks int
	Pages       map[string][]*Link
}

// Page is a struct used to store the details
// of an individual page of a website
type Page struct {
	URL    string
	Links  []Link
	Broken bool
}

var queue, broken map[string]bool
var pages map[string][]Link

var mu sync.RWMutex

func init() {
	queue = make(map[string]bool)
	broken = make(map[string]bool)
	pages = make(map[string][]Link)
}

// GenerateSiteMap will take in a url and
// crawls through all its pages and returns
// pages with all the urls corresponding to
// a page
func GenerateSiteMap(url string) {
	baseURL, err := getBaseURL(url)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	fmt.Println("Base URL :", baseURL)

	jobs := make(chan string, 100)
	// done := make(chan string, 100)

	jobs <- url

	time.Sleep(5 * time.Second)

	fmt.Println(len(pages), len(broken))
}

// GetBaseURL will return the base-url of
// the given url
func getBaseURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	baseURL := resp.Request.URL.Scheme + "://" + resp.Request.Host
	return baseURL, nil
}

func addToBroken(page Page) {
	mu.Lock()

	page.URL = strings.TrimRight(page.URL, "/")

	// Add page to the map of broken
	// and remove it from the queue
	broken[page.URL] = true
	delete(queue, page.URL)

	mu.Unlock()
}

// AddToQueue will add a url to the
// existing queue
func addToQueue(url string) {
	mu.Lock()

	url = strings.TrimRight(url, "/")
	queue[url] = true

	mu.Unlock()
}

// CleanURL will trim and clean the url to just
// the base url
func cleanURL(url string) string {
	cleanedURL := url

	containsHash := strings.ContainsAny(cleanedURL, "#")

	if containsHash {
		indx := strings.Index(cleanedURL, "#")
		cleanedURL = cleanedURL[:indx]
	}

	conatiainsQuery := strings.Contains(cleanedURL, "?")
	if conatiainsQuery {
		indx := strings.Index(cleanedURL, "?")
		cleanedURL = cleanedURL[:indx]
	}

	return cleanedURL
}

// recursive traditional way of crawling the
// web pages

// func recurParse(baseURL, url string) {
// 	links := ParseURL(baseURL, url)
// 	pages[url] = links

// 	// fmt.Printf("Links : %d\tPages : %d\tURL : %s\n", len(links), len(pages), url)
// 	for _, v := range links {
// 		// fmt.Println(v.URL)
// 		cleanedURL := cleanURL(v.URL)

// 		if strings.Contains(cleanedURL, baseURL) && pages[cleanedURL] == nil {
// 			recurParse(baseURL, cleanedURL)
// 			// addToQueue(baseURL, cleanedURL)
// 		}
// 	}
// }
