package main

import (
	"flag"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"
)

var concurrentNum int
var interval int

func init() {
	flag.IntVar(&concurrentNum, "cnum", 100, "默认并发数100")
	flag.IntVar(&interval, "interval", 0, "默认抓取间隔0秒")
}
func main() {
	flag.Parse()
	wg := sync.WaitGroup{}
	url := "https://apps.apple.com/cn/app/id414478124"
	wg.Add(concurrentNum)
	pool := make(chan int, 10)
	go func() {
		for i := 0; i < concurrentNum; i++ {
			pool <- i
		}
		close(pool)
	}()

	for taskId := range pool {
		time.Sleep(time.Duration(interval) * time.Second)
		go crawl(url, &wg, taskId, pool)
	}
	fmt.Println("wg is waiting")
	wg.Wait()
	fmt.Println("all finished")

}

func crawl(url string, wg *sync.WaitGroup, taskId int, pool chan int) {
	defer func() {
		err := recover()
		if err != nil {
			pool <- taskId
			switch err.(type) {
			case runtime.Error:
				fmt.Println("runtime error:", err)
			default:
				fmt.Println("error:", err)
			}
		}

	}()
	fmt.Println("当前抓取taskId为：", taskId)
	fetch(url, wg)
}
func fetch(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(url)

	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	parse(string(body))
}

func parse(content string) {
	doc, err := htmlquery.Parse(strings.NewReader(content))
	if err != nil {
		panic(err)
	}
	providerField := htmlquery.Find(doc, "//dt[@class='information-list__item__term medium-valign-top']/text()")[0]
	providerName := htmlquery.Find(doc, "//dd[@class='information-list__item__definition']//text()")[0]
	fmt.Println(strings.Trim(providerField.Data, " \n"), strings.Trim(providerName.Data, " \n"))
	sizeField := htmlquery.Find(doc, "//dt[@class='information-list__item__term medium-valign-top']/text()")[1]
	sizeName := htmlquery.Find(doc, "//dd[@class='information-list__item__definition']//text()")[1]
	fmt.Println(strings.Trim(sizeField.Data, " \n"), strings.Trim(sizeName.Data, " \n"))
	//fmt.Println(content)
}
