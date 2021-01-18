package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"sync"
	"os"
)

var wg sync.WaitGroup

func scnafile(file string, url chan string)  {
	files, err := os.Open(file)
	if err != nil {
		fmt.Println("open error")
	}
	defer files.Close()
	n := bufio.NewScanner(files)
	for n.Scan(){
		//fmt.Println(n.Text())
		url <- n.Text()
	}
	err_ := n.Err()
	if err_ != nil {
		fmt.Println("Error scanner")
	}
	close(url)

}

func reurl(url chan string)  {
	for{
		select {
		case xdd, ok := <-url:
			if !ok {
				wg.Done()
				return
			}
		//fmt.Println(xdd)
		c := colly.NewCollector()

		// Find and visit all links
		c.OnHTML("title", func(e *colly.HTMLElement) {
			fmt.Println(xdd, e.Text)
		})
		c.Visit(xdd)
		wg.Done()
		return
		}
	}
}

func main() {
	url :=make(chan string)
	var filename string
	flag.StringVar(&filename, "f", "", "请输入url文件名字")
	flag.Parse()
	if filename == ""{
		fmt.Println(`xxx`)
		os.Exit(1)
	}

	go scnafile(filename, url)
	for i := 0; i < 30; i++{
		wg.Add(1)
		go reurl(url)

	}
	wg.Wait()

}
