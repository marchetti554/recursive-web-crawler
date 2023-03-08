package main

import (
	"fmt"
	"testing"
	"time"
)

const testingWebSite = "https://quotes.toscrape.com/"

// In this test I would add some asserts regarding the directories, I didn't have more time to add that functionality
// that's why I'm just
// However, I found this site that seems pretty good to crawl, it's not too big, and it has some nested links.
func TestCrawl(t *testing.T) {
	start := time.Now()

	InitialURL = testingWebSite
	Crawl(testingWebSite)

	elapsed := time.Since(start)

	fmt.Println("Time crawling the web: ", elapsed)
}
