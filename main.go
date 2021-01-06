package main

import (
	"fmt"

	"github.com/ashishkhuraishy/sitemap_gen/htmlparser"
)

func main() {
	url := "https://pkg.go.dev/golang.org/x/net/html/"

	links := htmlparser.Parse(url)

	for _, v := range links {
		fmt.Println(v)
	}
}
