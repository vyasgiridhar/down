package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/urfave/cli"
)

func getHeadersAndStart(url, output string, g int) {

	client := new(http.Client)

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		fmt.Println("HEAD request failed")
		os.Exit(3)
	}

	header, httperr := client.Do(req)
	if httperr != nil {
		fmt.Println("HEAD request failed")
		os.Exit(3)
	}
	fmt.Println(header.Header)
	if header.Header.Get("Accept-Ranges") == "bytes" {
		fmt.Println("Downloading")
		length, err := strconv.Atoi(header.Header.Get("Content-Length"))
		if err != nil {
			fmt.Println("Content-Length could not be parsed")
			os.Exit(4)
		}
		if output == "" {
			output = strings.Replace(strings.Split(header.Header.Get("Content-Disposition"), "filename=")[1], "]", "", 2)
		}
		multiDownload(url, output, length, g)
	}

}

func downPart(wg *sync.WaitGroup, url string, dataChan chan []byte, range1, range2 int) {

	defer wg.Done()
	client := new(http.Client)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("GET request failed")
		os.Exit(2)
	}

	req.Header.Set("Range", strconv.Itoa(range1)+"-"+strconv.Itoa(range2))
	data, err := client.Do(req)
	if err != nil {
		fmt.Println("GET request failed")
		os.Exit(2)
	}
	dataByte := new(bytes.Buffer)
	dataByte.ReadFrom(data.Body)

	dataChan <- dataByte.Bytes()
}

func multiDownload(url, output string, length, g int) bool {
	var wg sync.WaitGroup

	split := length / g
	fmt.Println(length)
	dataChan := make(chan []byte)

	for i := 0; i < g; i++ {
		wg.Add(1)

		range1 := i * split
		range2 := (i+1)*split - 1
		if range2 == length-2 {
			range2 = length
		}
		go downPart(&wg, url, dataChan, range1, range2)
	}

	go func() {
		wg.Wait()
		close(dataChan)
	}()

	for data := range dataChan {
		if output == "" {
			fileName := strings.Split(url, "/")
			output = fileName[len(fileName)-1]
		}
		ioutil.WriteFile(output, data, os.ModeAppend)
		os.Chmod(output, 0777)
	}
	return true
}

func main() {

	app := cli.NewApp()
	app.Name = "down"
	app.Usage = "Multigoroutine downloader"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "threads, g",
			Value: 4,
			Usage: "Number of simultaneous threads to use",
		},
		cli.StringFlag{
			Name:  "output, o",
			Value: "",
			Usage: "Output file name",
		},
	}

	app.Action = func(c *cli.Context) error {
		var url string
		if c.NArg() > 0 {
			url = c.Args().Get(0)
		}
		if url == "" {
			fmt.Println("URL not provided")
			os.Exit(1)
		} else {
			g := c.Int("threads")
			output := c.String("output")
			if g == 0 {
				g = 1
			}
			getHeadersAndStart(url, output, g)
		}
		return nil
	}

	app.Run(os.Args)

}
