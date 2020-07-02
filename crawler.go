package crawler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// CrawlBot defines the control for running the crawling
type CrawlBot struct {
	handler     Handler
	urlsCrawled []string
	in          chan Request
	out         chan IResult
	wg          sync.WaitGroup
	client      *http.Client
	doneTimer   *time.Timer // stop crawling when in channel is idle for a specified time
}

// NewCrawlBot returns an initialized CrawBot and starts the process goroutine
func NewCrawlBot(h Handler) *CrawlBot {
	bot := &CrawlBot{
		handler:     h,
		urlsCrawled: []string{},
		in:          make(chan Request, 1),
		out:         make(chan IResult, 1),
		client:      http.DefaultClient,
	}

	// start the done timer
	bot.doneTimer = time.NewTimer(100 * time.Millisecond)

	// start the go routine to process request
	bot.wg.Add(2)
	go bot.processRequests()
	go bot.processResults()

	return bot
}

// Crawl accepts new requests and process them
func (c *CrawlBot) Crawl(method string, urls ...string) {
	// for every new request, send to the in channel
	for _, u := range urls {
		parsedURL, err := url.Parse(u)
		if err != nil {
			log.Fatalf("URL %s can not be parsed, moving on", u)
			continue
		}
		c.in <- Request{parsedURL, method}
	}
}

// Done : Close request channel and wait for goroutine to finish
func (c *CrawlBot) Done() {
	//
	// close(c.in)
	// close(c.out)
	c.wg.Wait()
}

// Process Requests in a goroutine, if in channel is idle for a certain time
// close in channel
func (c *CrawlBot) processRequests() {
	// defer c.wg.Done()

processLoop:
	for {
		select {
		case r := <-c.in:
			// New Request, reset timer
			c.doneTimer.Stop()
			c.doneTimer.Reset(100 * time.Millisecond)
			fmt.Printf("timer reset, got request %s\n", r.URL.String())
			//
			go func(re Request) {
				c.out <- c.doRequest(re)
			}(r)
		case <-c.doneTimer.C:
			// channel idle for enough time, stop crawler
			close(c.in)
			c.out <- nil // send nil to out as close signal
			c.wg.Done()
			break processLoop
		}
	}

	// for r := range c.in {
	// 	// create network request go routine
	// 	go func(r Request) {
	// 		c.out <- c.doRequest(r)
	// 	}(r)
	// }
}

// processs the results from make request goroutine
func (c *CrawlBot) processResults() {
	for res := range c.out {
		// fmt.Printf("Got result for %s\n", res.Request().URL.String())
		if res == nil {
			// close signal, could still be result to process
			close(c.out)
			c.wg.Done()
			break
		}
		c.handler.Handle(res.Request(), res.Response(), res.Error())
		c.urlsCrawled = append(c.urlsCrawled, res.Request().URL.String())
	}
}

// Process individual request
func (c *CrawlBot) doRequest(r Request) IResult {
	req, err := http.NewRequest(r.method, r.URL.String(), nil)
	if err != nil {
		return &Result{r, nil, err}
	}

	res, err := c.client.Do(req)
	if err != nil {
		return &Result{r, nil, err}
	}
	return &Result{r, res, nil}
}
