package util

import (
	"golang.org/x/net/html"
	"io"
	"sync"
)

func GetHrefLinks(page io.Reader) (*sync.Map, error) {
	tokenizer := html.NewTokenizer(page)
	var hrefLinksMap = &sync.Map{}
	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			return hrefLinksMap, nil
		}
		token := tokenizer.Token()
		if tt == html.StartTagToken && token.DataAtom.String() == "a" {
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					hrefLinksMap.Store(attr.Val, struct{}{})
				}
			}
		}
	}
}
