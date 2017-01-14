package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
)

func getHeaders(url string) (header *http.Response) {
	client := new(http.Client)

	req, _ := http.NewRequest("HEAD", url, nil)
	header, _ = client.Do(req)
	return
}

func downPart(url string, dataChan chan []byte, range1, range2 int) []byte {
	client := new(http.Client)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Range", strconv.Itoa(range1)+"-"+strconv.Itoa(range2))
	data, _ := client.Do(req)
	dataByte := new(bytes.Buffer)
	dataByte.ReadFrom(data.Body)
	fmt.Println("Done")
	return dataByte.Bytes()
}
func multiDownload(url string, length int) bool {
	x := 4
	split := length / x
	fmt.Println(length)
	dataChan := make([]chan []byte, x)
	for i := 0; i < x; i++ {
		range1 := i * split
		range2 := (i+1)*split - 1
		if range2 == length-2 {
			range2 = length
		}
		fmt.Println(len(dataChan), range1, range2)
		go downPart(url, dataChan[i], i*split, i+1*split)
	}
	//ioutil.WriteFile("ba.jpg", dataByte.Bytes(), 0777)
	return true
}

func main() {
	header := getHeaders("https://i.redd.it/sjqcrazacd9y.jpg")
	if header.Header.Get("Accept-Ranges") == "bytes" {
		fmt.Println("Downloading")
		length, _ := strconv.Atoi(header.Header.Get("Content-Length"))
		multiDownload("https://i.redd.it/sjqcrazacd9y.jpg", length)
	}

}
