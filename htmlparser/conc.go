package htmlparser

import (
	"fmt"
	"strings"
	"sync"
)

var wg, jp sync.WaitGroup

// Create worker
//	[] will listen endlessly to jobs
//	[] process each job
//	[] spawn a resultpool
//	[] trigger done channel once the job is finished
func worker(baseURL string, jobs <-chan string, done chan<- string) {
	for job := range jobs {
		// Spawn a goroutine to process the resullts
		// and add them to queue
		wg.Add(1)
		page := ParseURL(baseURL, job)
		go sheduler(baseURL, page)
		done <- job
	}
}

// Sheduler
//	[] if the result has an error add it to broken
//	[] loop through each result
//	[] clean the url
//	[] if the url has baseUrl & url is not
//	   present in the queue then add that
//	   to the queue
func sheduler(baseURL string, page Page) {
	if page.Broken {
		// TODO: add page to broken pages
		addToBroken(page)
		return
	}

	// TODO: add page url and its links
	// to the result
	addToPages(page)

	// validLinks := make([]string, 0)

	for _, link := range page.Links {
		// TODO: check if the links are already
		// processed / in the queue. if not add
		// them to the queue
		url := cleanURL(link.URL)

		if strings.Contains(url, baseURL) && !existingURL(url) {
			// validLinks = append(validLinks, url)
			addToQueue(url)
		}
	}

	wg.Done()
}

// Done Loop
// when a job finishes its task that job is said to
// be done. When a done channel is triggered then
// then take the next element from queue and assign
// it to a job
func assignjobs(done <-chan string, jobs chan<- string) {
	for range done {
		// TODO: check if the queue if empty
		if job := getElementFromQueue(); job != "" {
			jobs <- job
		} else {
			// Close all the channels
			return
		}

	}
}

func getElementFromQueue() string {
	mu.Lock()
	defer mu.Unlock()
	if len(queue) == 0 {
		return ""
	}

	for url := range queue {
		delete(queue, url)
		return url
	}

	return ""
}

func addToPages(page Page) {
	mu.Lock()

	// Add page to the map of pages
	// and remove it from the queue
	pages[page.URL] = page.Links

	fmt.Printf("Pages : %d\t Queue : %d\n", len(pages), len(queue))
	// fmt.Println(queue)

	mu.Unlock()
}

// CheckUrl will check the url is already processed
// or already in te queue
func existingURL(url string) bool {
	mu.Lock()
	defer mu.Unlock()

	// fmt.Println(pages[url])

	// Make all the urls in the same format
	url = strings.TrimRight(url, "/")

	// Check if the page is already processed
	if len(pages[url]) == 0 {
		// fmt.Println(pages[url])
		// Check if it is in the queue
		if !queue[url] && !broken[url] {
			return false
		}

	}

	return true
}

func jobpool(baseURL string, jobs chan string, results chan Page, quit chan int) {
	// mark the jobpool as done
	defer jp.Done()

	sem := make(chan int, 100)

	for {
		select {
		case job := <-jobs:
			sem <- 1
			wg.Add(1)
			go getpage(baseURL, job, results, sem)
			// case result := <-results:
			// sem <- 1
			wg.Add(1)
			// go sheduler(baseURL, result, jobs, sem)
		case _, ok := <-quit:
			if !ok {
				fmt.Println("Quitting")
				return
			}

		}

	}
}

func getpage(baseURL, url string, result chan Page, sem chan int) {
	// fmt.Println("Adding", url, "to the queue")
	addToQueue(url)

	result <- ParseURL(baseURL, url)
	<-sem

	// fmt.Println("Done Processing", url)
	wg.Done()
}
