package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func getHeaders(url string) (header *http.Response) {
	client := new(http.Client)

	req, _ := http.NewRequest("HEAD", url, nil)
	header, _ = client.Do(req)
	return
}

func multiDownload(url string, length int) bool {
	//x := 4
	//split := length / x
	client := new(http.Client)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Range", "0-"+strconv.Itoa(length))
	data, _ := client.Do(req)
	dataByte := new(bytes.Buffer)
	dataByte.ReadFrom(data.Body)
	ioutil.WriteFile("ba.jpg", dataByte.Bytes(), 0777)
	return true
}

func main() {
	header := getHeaders("https://dl.dropboxusercontent.com/content_link/b4uW3v4f4qgYyPsKbPDBdKYk529ukI9edxr3IhNQt8GXeyghl3ZPvLcubkfJvtai/file?dl=1")
	if header.Header.Get("Accept-Ranges") == "bytes" {
		fmt.Println("Downloading")
		length, _ := strconv.Atoi(header.Header.Get("Content-Length"))
		multiDownload("https://dl.dropboxusercontent.com/content_link/b4uW3v4f4qgYyPsKbPDBdKYk529ukI9edxr3IhNQt8GXeyghl3ZPvLcubkfJvtai/file?dl=1", length)
	}

}
