package crawler

import (
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
	out         chan Request
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
		out:         make(chan Request, 1),
		client:      http.DefaultClient,
	}

	// start the done timer, arbitrary time for first timer
	bot.doneTimer = time.NewTimer(1000 * time.Millisecond)

	// start the go routine to process request
	bot.wg.Add(1)
	go bot.processRequests()
	go bot.logResults()

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
processLoop:
	for {
		select {
		case r := <-c.in:
			// New request, stop timer
			go c.stopTimer()
			//
			go func(re Request) {
				res := c.doRequest(re)
				c.handler.Handle(res.Request(), res.Response(), res.Error())
				c.out <- re
				// Request finished, reset timer
				c.doneTimer.Reset(10 * time.Millisecond)
			}(r)
		case <-c.doneTimer.C:
			// channel idle for enough time, stop crawler
			close(c.in)
			close(c.out)
			c.wg.Done()
			break processLoop
		}
	}
}

// run continuous stop to make sure the timer is stopped and resetted
// on the last request
func (c *CrawlBot) stopTimer() {
	for {
		if c.doneTimer.Stop() {
			break
		}
	}
}

// log the results when request is finished handling
func (c *CrawlBot) logResults() {
	for res := range c.out {
		c.urlsCrawled = append(c.urlsCrawled, res.URL.String())
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
