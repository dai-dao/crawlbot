package crawler

import (
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestFetch(t *testing.T) {
	t.Run("test run HEAD request", func(t *testing.T) {
		fmt.Print("\nTest run HEAD request\n")
		// Initialize and Run the crawler
		crawler := NewCrawlBot(handler)
		crawler.Crawl("HEAD", "http://golang.org", "http://google.com")
		crawler.Done()
		// assert
		want := []string{"http://golang.org", "http://google.com"}
		if len(crawler.urlsCrawled) != 2 {
			t.Fatalf("Wanted %v, got %v", want, crawler.urlsCrawled)
		}
	})

	t.Run("test when crawler is done, can not crawl anymore and return error", func(t *testing.T) {
		fmt.Print("\nTest when crawler is done, can not crawl anymore and return error\n")
		crawler := NewCrawlBot(handler)
		err1 := crawler.Crawl("HEAD", "http://golang.org", "http://google.com")
		crawler.Done()
		err2 := crawler.Crawl("HEAD", "http://golang.org", "http://google.com")
		//
		if err1 != nil {
			log.Fatalf("Expected nil error, got %v", err1)
		}
		if err2 != ErrStopped {
			log.Fatalf("Expected error %s, got %v", ErrStopped, err2)
		}
		//
		want := []string{"http://golang.org", "http://google.com"}
		if len(crawler.urlsCrawled) != 2 {
			t.Fatalf("Wanted %v, got %v", want, crawler.urlsCrawled)
		}
	})

	t.Run("test crawler can be stopped at any time without waiting to finish", func(t *testing.T) {
		fmt.Print("\nTest crawler can be stopped at any time without waiting to finish\n")
		crawler := NewCrawlBot(handler)
		crawler.Crawl("HEAD", "http://golang.org", "http://google.com")
		// Stop crawling goroutine after some time or condition
		timer := time.NewTimer(10 * time.Millisecond)
		go func() {
			<-timer.C
			crawler.Stop()
		}()
		//
		crawler.Done()
		//
		if len(crawler.urlsCrawled) == 2 {
			t.Fatalf("should be less than 2, but got 2")
		}
	})
}

var handler = HandlerFunc(func(r Request, res *http.Response, err error) {
	if err != nil {
		log.Fatalf("Error doing request with url %s, %s", r.URL.String(), err)
	}
	fmt.Printf("[%d] %s %s\n", res.StatusCode, r.method, r.URL.String())
})
