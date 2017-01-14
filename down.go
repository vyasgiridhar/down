package main

import (
	"fmt"
	"net/http"
)

func getHeaders(url string) (header *http.Response) {
	client := new(http.Client)

	req, _ := http.NewRequest("HEAD", url, nil)
	header, _ = client.Do(req)
	return
}

func main() {
	header := getHeaders("https://i.redd.it/uv1x1l3gtg9y.jpg")
	if header.Header.Get("Accept-Ranges") == "bytes" {
		fmt.Println("Downloading")

	}
}
