package crawler

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestFetch(t *testing.T) {
	t.Run("test run HEAD request", func(t *testing.T) {
		// Initialize and Run the crawler
		crawler := NewCrawlBot(handler)

		// new links will be added continuously, and asynchronously
		crawler.Crawl("HEAD", "http://golang.org", "http://google.com")

		// wait for crawl
		// time.Sleep(200 * time.Millisecond)

		// done
		crawler.Done()

		// assert
		want := []string{"http://golang.org", "http://google.com"}
		if len(crawler.urlsCrawled) != 2 {
			t.Fatalf("Wanted %v, got %v", want, crawler.urlsCrawled)
		}
	})

	t.Run("test when crawler is done, can not crawl anymore and return error", func(t *testing.T) {

	})

	t.Run("test crawler can be stopped at any time without waiting to finish", func(t *testing.T) {

	})
}

var handler = HandlerFunc(func(r Request, res *http.Response, err error) {
	if err != nil {
		log.Fatalf("Error doing request with url %s, %s", r.URL.String(), err)
	}
	fmt.Printf("[%d] %s %s\n", res.StatusCode, r.method, r.URL.String())
})
