package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sync"
)

/*

用的colly模块，-f指定文件名

*/
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
		client := &http.Client{}
		url := xdd
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}
		resp, err := client.Do(req)
		defer resp.Body.Close()
		content, _ := ioutil.ReadAll(resp.Body)
		r := regexp.MustCompile(`<title>(.*?)</title>`)
		ret := r.FindAllStringSubmatch(string(content),-1)
		for _, substr := range ret  {
			fmt.Printf("[%d+] %s %s %s \n",resp.StatusCode, substr[1], xdd, resp.Header.Get("Server"))
		}
		wg.Done()
		return
		}
	}
}

func main() {
	url :=make(chan string)
	var filename string

	flag.StringVar(&filename, "f", "1.txt", "请输入url文件名字")
	flag.Parse()
	if filename == ""{
		fmt.Println(`Mylogo`)
		os.Exit(1)
	}

	go scnafile(filename, url)
	for i := 0; i < 30; i++{
		wg.Add(1)
		go reurl(url)
	}
	wg.Wait()

}
