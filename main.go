package main

import (
	"github.com/ashishkhuraishy/sitemap_gen/htmlparser"
)

func main() {
	// url := "https://pkg.go.dev/golang.org/x/net/html/"
	url := "https://go.dev/"

	// links := htmlparser.ParseURL(url)
	htmlparser.GenerateSiteMap(url)

	// 	for _, v := range links {
	// 		fmt.Println(v)
	// 	}
}
