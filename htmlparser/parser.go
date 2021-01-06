package htmlparser

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// Link is a struct which will store
// the node which contains the link,
// the url and the text inside the
// `a` tag
type Link struct {
	URL  string
	Text string
}

// Parse will take in a url
// and parse it to return all
// the links available on the
// html page
func Parse(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer resp.Body.Close()

	html, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	nodes := linkNodes(html)

	if last := len(url) - 1; last >= 0 && url[last] == '/' {
		url = url[:last]
		fmt.Println(url)
	}

	var links []Link
	for _, n := range nodes {
		links = append(links, getLink(n, url))
	}

	for _, v := range links {
		fmt.Println(v.URL)
	}

}

// linkNodes will recursively loop through all the
// elements in the node and return all the nodes
// containing the `a` tag
func linkNodes(node *html.Node) []*html.Node {
	if node.Type == html.ElementNode && node.Data == "a" {
		return []*html.Node{node}
	}

	var nodes []*html.Node
	for i := node.FirstChild; i != nil; i = i.NextSibling {
		nodes = append(nodes, linkNodes(i)...)
	}

	return nodes
}

// getLink fn will take a node and
// convert it into struct of link
func getLink(node *html.Node, url string) Link {
	var link Link

	link.Text = getText(node)
	for _, n := range node.Attr {
		if n.Key == "href" && len(n.Val) > 0 {
			// fmt.Println(n.Val)
			link.URL = n.Val
			if n.Val[0] == '/' || n.Val[0] == '#' {
				link.URL = url + n.Val
			}

			return link
		}
	}

	return link
}

// GetText will loop through all the elemnts
// on a given node, combines and returns the
// text from all the nodes
func getText(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}

	if node.Type != html.ElementNode {
		return ""
	}

	var text string
	for i := node.FirstChild; i != nil; i = i.NextSibling {
		text += strings.TrimSpace(getText(i)) + " "
	}

	// text = strings.Join(strings.Fields(text), " ")

	return text
}
