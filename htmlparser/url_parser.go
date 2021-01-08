package htmlparser

import (
	"log"
	"net/http"
	"strings"
	"time"

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

// ParseURL will take in a url
// and parse it to return all
// the links available on the
// html page
func ParseURL(baseURL, url string) Page {
	url = strings.TrimRight(url, "/")
	// fmt.Println("Running", url)
	page := Page{
		URL: url,
	}

	resp, err := http.Get(url)
	if err != nil {
		if strings.Contains(err.Error(), "too many open") {
			time.Sleep(2 * time.Second)
			return ParseURL(baseURL, url)
		}
		log.Println(err.Error())
		// errchan <- page
		page.Broken = true
		return page
	}
	defer resp.Body.Close()

	html, err := html.Parse(resp.Body)
	if err != nil {
		log.Println(err.Error())
		// errchan <- page
		page.Broken = true
		return page
	}

	nodes := linkNodes(html)

	var links []Link
	for _, node := range nodes {
		links = append(links, getLink(node, baseURL, url))
	}

	// page.URL = url
	page.Links = links

	return page
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
func getLink(node *html.Node, baseURL, url string) Link {
	var link Link

	link.Text = getText(node)
	for _, n := range node.Attr {
		if n.Key == "href" && len(n.Val) > 0 {
			// fmt.Println(n.Val)
			link.URL = n.Val
			if n.Val[0] == '/' || n.Val[0] == '#' || n.Val[0] == '?' {
				link.URL = baseURL + n.Val
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
		text += getText(i)
	}

	// converting the text into []string without space
	// then joining each other with one white-space in
	// between. (HTML pages can have uneven whitespaces
	// which can become messy sometimes)
	text = strings.Join(strings.Fields(text), " ")

	return text
}
