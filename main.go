package main

import (
	"sync"

	"github.com/ashishkhuraishy/sitemap_gen/htmlparser"
)

var wg, jwg sync.WaitGroup

var mp sync.Map

func main() {
	url := "https://pkg.go.dev/golang.org/x/net/html/"
	// url := "https://go.dev/"
	// url := "https://www.vervesearch.com/blog/how-to-make-a-simple-web-crawler-in-go/"

	// links := htmlparser.ParseURL(url)
	htmlparser.GenerateSiteMap(url)

}
