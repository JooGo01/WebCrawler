package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/steelx/extractlinks"
)

var (
	config = &tls.Config{
		InsecureSkipVerify: true,
	}

	transport = &http.Transport{
		TLSClientConfig: config,
	}

	netClient = &http.Client{
		Transport: transport,
	}

	queue = make(chan string)

	visit = make(map[string]bool)
)

func main() {

	arguments := os.Args[1:]

	if len(arguments) == 0 {
		fmt.Println("Missing URL, e.g. go-webscrapper http://js.org/")
		os.Exit(1)
	}

	go func() {
		queue <- arguments[0]
	}()

	for href := range queue {
		if !visit[href] {
			crawlURL(href)
		}
	}

}

func crawlURL(href string) {

	visit[href] = true

	fmt.Printf("Crawling the URL %v \n", href)

	response, err := netClient.Get(href)

	checkError(err)

	defer response.Body.Close()

	links, err := extractlinks.All(response.Body)

	checkError(err)

	for _, link := range links {
		absURL := fixURL(link.Href, href)
		go func() {
			queue <- absURL
		}()
	}

}

func fixURL(href, URL string) string {

	uri, err := url.Parse(href)

	if err != nil {
		return ""
	}

	base, err := url.Parse(URL)
	if err != nil {
		return ""
	}

	uriFix := base.ResolveReference(uri)

	return uriFix.String()
}

func checkError(err error) {

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
