package htmlparser

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// SiteMap is a struct which stores all
// the info about the website
type SiteMap struct {
	Domain      string
	PageCount   int
	BrokenLinks int
	Pages       map[string][]*Link
}

var queue map[string]bool
var pages map[string][]*Link

func init() {
	queue = make(map[string]bool)
	pages = make(map[string][]*Link)
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

	recurParse(baseURL, url)

	sitemap := SiteMap{
		Domain:      baseURL,
		PageCount:   len(pages),
		BrokenLinks: 0,
		Pages:       pages,
	}

	fmt.Println(sitemap.Domain, sitemap.PageCount)

	for page, links := range pages {
		fmt.Println(page)
		for _, link := range links {
			fmt.Println("\t", link)
		}
	}
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

func recurParse(baseURL, url string) {
	links := ParseURL(baseURL, url)
	pages[url] = links

	// fmt.Printf("Links : %d\tPages : %d\tURL : %s\n", len(links), len(pages), url)
	for _, v := range links {
		// fmt.Println(v.URL)
		cleanedURL := cleanURL(v.URL)

		if strings.Contains(cleanedURL, baseURL) && pages[cleanedURL] == nil {
			recurParse(baseURL, cleanedURL)
			// addToQueue(baseURL, cleanedURL)
		}
	}
}

// AddToQueue will add a url to the
// existing queue
func addToQueue(baseURL, url string) {
	if queue[url] {
		return
	}
	fmt.Printf("Queue %d\tURL : %s\n", len(queue), url)

	queue[url] = true
	recurParse(baseURL, url)
	delete(queue, url)
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
