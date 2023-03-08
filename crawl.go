package main

import (
	"crawler/util"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
)

var (
	visitedLinks = &sync.Map{}
	InitialURL   string
)

/*
1. Check first webpage for links
2. Fix each found link adding baseURL when necessary
3. Save each page in a different map in order to be able to track visited links
4. Check each link for the map and do the same
*/

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("... Exiting program ...")
			os.Exit(1)
		}
	}()

	//Get argument
	flag.Parse()
	webPage := flag.Args()

	//Check argument
	if len(webPage) != 1 {
		fmt.Println("... Only one argument is required ...") //Add directory as 2nd argument
		os.Exit(1)
	}

	InitialURL = webPage[0]

	//Lets crawl!
	Crawl(InitialURL)

	fmt.Println("... Finished crawling ...")
}

func Crawl(webPage string) {
	queue := make(chan string)

	go func() { queue <- webPage }()

	visitedLinks.Store(webPage, struct{}{})

	done := make(chan bool)

	for i := 0; i < 3; i++ {
		go func() {
			for uri := range queue {
				addToQueue(uri, queue)
			}
			done <- true
		}()
	}
	<-done
}

func addToQueue(url string, queue chan string) {
	fmt.Println("getting web ---> ", url)
	web, err := util.GetWeb(url)
	if err != nil {
		fmt.Printf("%+v", err)

		return
	}

	links, err := util.GetHrefLinks(web)
	if err != nil {
		fmt.Printf("%+v", err)

		return
	}

	links.Range(func(link, _ any) bool {
		fixedUrl := fixUrl(fmt.Sprintf("%s", link))
		if fixedUrl == "" {
			return true
		}

		_, ok := visitedLinks.Load(fixedUrl)
		if !ok {
			go func() {
				Crawl(fixedUrl)
			}()
			go func() {
				queue <- fixedUrl
			}()
		}
		return true
	})

}

func fixUrl(href string) string {
	linkUrl, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseUrl, err := linkUrl.Parse(InitialURL)
	if err != nil {
		return ""
	}
	linkUrl = baseUrl.ResolveReference(linkUrl)
	if strings.HasPrefix(linkUrl.String(), InitialURL) {
		return linkUrl.String()
	}
	return ""
}
