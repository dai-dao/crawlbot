package main

type CrawlBot struct {
	h func()

	urlsCrawled []string

	in chan string
}

func NewCrawlBot(h func()) CrawlBot {
	bot := CrawlBot{h, []string{}, make(chan string)}

	// start the go routine to process request
	go processRequests()

	return bot
}

func (c *CrawlBot) Crawl(method string, urls ...string) {
	// for every new request, send to the in channel
	for _, u := range urls {
		c.in <- u
	}

	c.urlsCrawled = append(c.urlsCrawled, u)

	c.h()
}

func (c *CrawlBot) Done() {

}

// Private internal function
func (c *CrawlBot) processRequests() {
	for u := range c.in {

	}
}
