package main

import (
	"bufio"
	"crypto/tls"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
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
	//n := bufio.NewReader(files)
	//for {
	//	a, _, c := n.ReadLine()
	//	if c == io.EOF {
	//		break
	//	}
	//	fmt.Println(string(a),"aaa")
	//	url <- string(a)
	//}
	n := bufio.NewScanner(files)
	for n.Scan(){
		fmt.Println(n.Text())
		url <- n.Text()
		continue
	}
	err_ := n.Err()
	if err_ != nil {
		fmt.Println("Error scanner")
	}
	close(url)

}

func reurl(url chan string)  {
	timeUnix := time.Now().Format("2006-01-02 15:04:05")+".csv"
	for{
		select {
		case xdd, ok := <-url:
			if !ok {
				wg.Done()
				return
			}
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify:true},
		}
		client := &http.Client{Transport:tr}
		url := xdd
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("url2 error")
			continue
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("url: %s error\n", url)
			continue
		}
		defer resp.Body.Close()

		content, _ := ioutil.ReadAll(resp.Body)
		r := regexp.MustCompile(`<title>(.*)</title>`)
		ret := r.FindAllStringSubmatch(string(content),-1)
		//fmt.Println(ret)
		for _, substr := range ret  {
			fmt.Printf("[%d+] %s %s %s \n",resp.StatusCode, substr[1], xdd, resp.Header.Get("Server"))
			writeurl(timeUnix, strconv.Itoa(resp.StatusCode), xdd, substr[1], resp.Header.Get("Server"))
		}
		case <- time.After(time.Duration(2) * time.Second):
			wg.Done()
			return


		}
	}
}

//写文件
func writeurl(time_ string, code string, url string, title string, serverver string)  {
	file , err := os.OpenFile(time_, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	_, _ = file.WriteString("\xEF\xBB\xBF")
	if err != nil {
		fmt.Println("open file err:", err)
	}
	defer file.Close()
	w := csv.NewWriter(file)
	_ = w.Write([]string{url, title, code, serverver})
	w.Flush()
}

func main() {
	url :=make(chan string)
	var filename string

	flag.StringVar(&filename, "f", "", "请输入url文件名字")
	flag.Parse()
	if filename == ""{
		fmt.Println(`
 _____  _  _____  _    ____    ___   ___   __   _   _ 
(_   _)( )(_   _)( )  (  __)  /  _) (  _) (  ) ( \ ( )
  | |  | |  | |  | |  | |_    \_"-. | |   /o \ | \\| |
  ( )  ( )  ( )  ( )_ (  _)    __) )( )_ ( __ )( )\\ )
  /_\  /_\  /_\  /___\/____\  /___/ /___\/_\/_\/_\ \_\
                                               v1.0
`)
		os.Exit(1)
	}

	go scnafile(filename, url)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go reurl(url)
	}
	wg.Wait()

}
