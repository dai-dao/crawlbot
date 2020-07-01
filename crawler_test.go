package main

import (
	"testing"
)

func TestFetch(t *testing.T) {
	t.Run("test run HEAD request", func(t *testing.T) {
		// Initialize and Run the crawler
		crawler := NewCrawlBot(handler)

		// new links will be added continuously, and asynchronously
		crawler.Crawl("HEAD", "http://golang.org", "http://google.com")

		// done
		crawler.Done()

		// assert
		if len(crawler.urlsCrawled) != 2 {
			t.Errorf("num url crawled is wrong, got %d but want %d", len(crawler.urlsCrawled), 2)
		}
	})
}

func handler() {
}
