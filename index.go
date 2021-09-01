package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/steelx/extractlinks"
)

func main() {
	URL := "https://www.youtube.com/watch?v=2wmkHFTaXfA&t=767s"

	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	transport := &http.Transport{
		TLSClientConfig: config,
	}

	netClient := &http.Client{
		Transport: transport,
	}

	response, err := netClient.Get(URL)
	checkError(err)
	defer response.Body.Close()

	links, err := extractlinks.All(response.Body)

	fmt.Println(links)

}

func checkError(err error) {

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
