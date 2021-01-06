package main

import (
	"io"
	"log"
	"net/http"

	"github.com/ashishkhuraishy/sitemap_gen/htmlparser"
)

func main() {
	url := "https://pkg.go.dev/golang.org/x/net/html/"

	htmlparser.Parse(url)
}

func getHTML(url string) io.Reader {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	return resp.Body
}
