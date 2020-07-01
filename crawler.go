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
	out         chan Result
	wg          sync.WaitGroup
	client      *http.Client
}

// NewCrawlBot returns an initialized CrawBot and starts the process goroutine
func NewCrawlBot(h Handler) *CrawlBot {
	bot := &CrawlBot{
		handler:     h,
		urlsCrawled: []string{},
		in:          make(chan Request),
		out:         make(chan Result),
		client:      http.DefaultClient,
	}

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
	time.Sleep(200 * time.Millisecond)
	//
	close(c.in)
	close(c.out)
	c.wg.Wait()
}

// Private internal function to process Requests in a goroutine
func (c *CrawlBot) processRequests() {
	defer c.wg.Done()

	for r := range c.in {
		// create network request go routine
		go func(r Request) {
			c.out <- c.doRequest(r)
		}(r)
	}
}

// processs the results from make request goroutine
func (c *CrawlBot) processResults() {
	defer c.wg.Done()

	for res := range c.out {
		c.handler.Handle(res.req, res.res, res.err)
		c.urlsCrawled = append(c.urlsCrawled, res.req.URL.String())
	}
}

// Process individual request
func (c *CrawlBot) doRequest(r Request) Result {
	req, err := http.NewRequest(r.method, r.URL.String(), nil)
	if err != nil {
		return Result{r, nil, err}
	}

	res, err := c.client.Do(req)
	if err != nil {
		return Result{r, nil, err}
	}
	return Result{r, res, nil}
}
