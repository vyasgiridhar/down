package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func getHeaders(url string) (header *http.Response) {
	client := new(http.Client)

	req, _ := http.NewRequest("HEAD", url, nil)
	header, _ = client.Do(req)
	return
}

func downPart(wg *sync.WaitGroup, url string, dataChan chan []byte, range1, range2 int) {

	defer wg.Done()
	client := new(http.Client)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Range", strconv.Itoa(range1)+"-"+strconv.Itoa(range2))
	data, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	dataByte := new(bytes.Buffer)
	dataByte.ReadFrom(data.Body)
	fmt.Println("Done")

	dataChan <- dataByte.Bytes()
}
func multiDownload(url string, length int) bool {
	var wg sync.WaitGroup

	x := 4
	split := length / x
	fmt.Println(length)
	dataChan := make([]chan []byte, x)
	for i := range dataChan {
		dataChan[i] = make(chan []byte, 1)
	}

	for i := 0; i < x; i++ {
		wg.Add(1)

		range1 := i * split
		range2 := (i+1)*split - 1
		if range2 == length-2 {
			range2 = length
		}
		fmt.Println(len(dataChan), range1, range2)
		go downPart(&wg, url, dataChan[i], range1, range2)
	}
	fmt.Println("waiting")

	var data []byte
	wg.Wait()

	for i := 0; i < x; i++ {
		fmt.Println("waiting")
		data = <-dataChan[i]
		if i == 0 {
			ioutil.WriteFile("b.jpg", data, 0777)
		} else {
			ioutil.WriteFile("b.jpg", data, os.ModeAppend)
		}
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
