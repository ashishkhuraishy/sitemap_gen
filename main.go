package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/ashishkhuraishy/sitemap_gen/htmlparser"
)

var wg, jwg sync.WaitGroup

var mp sync.Map

func main() {
	// url := "https://pkg.go.dev/golang.org/x/net/html/"
	url := "https://go.dev/"

	// links := htmlparser.ParseURL(url)
	htmlparser.GenerateSiteMap(url)

}

// Test Concurrency code

func cunc(links []htmlparser.Link) {
	for _, v := range links {
		fmt.Println(v)
	}

	quit := make(chan int)
	ch := make(chan int)
	sem := make(chan int, 100)

	jwg.Add(1)
	go jobPool(ch, quit)

	for i := 0; i < 100000; i++ {
		sem <- 1
		go worker(ch, sem, i+1)
		wg.Add(1)

	}

	wg.Wait()
	close(quit)

	jwg.Wait()
}

func worker(ch, sem chan int, idx int) {
	rand.Seed(time.Now().Unix())
	// fmt.Println("Running channel", idx)
	time.Sleep(time.Duration(rand.Intn(3)) * time.Second)

	ch <- idx
	// fmt.Println("Finished", idx)
	<-sem
	wg.Done()
}

func jobPool(ch chan int, quit chan int) {
	fmt.Println("Starting pool")

	done := false

	for !done {
		select {
		case n := <-ch:
			fmt.Println("Recived", n)
		case _, ok := <-quit:
			if !ok {
				done = true
				break
			}
		}
	}

	fmt.Println("Exiting job pool")
	jwg.Done()
}
