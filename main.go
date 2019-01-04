package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sclevine/agouti"
)

func main() {
	url := "http://news.livedoor.com/topics/category/main/"
	driver := agouti.ChromeDriver()

	err := driver.Start()
	if err != nil {
		log.Printf("Failed to start driver: %v", err)
	}
	defer driver.Stop()

	page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		log.Printf("Failed to open page: %v", err)
	}

	err = page.Navigate(url)
	if err != nil {
		log.Printf("Failed to navigate: %v", err)
	}

	curContentsDom, err := page.HTML()
	if err != nil {
		log.Printf("Failed to get html: %v", err)
	}

	readerCurContents := strings.NewReader(curContentsDom)
	contentsDom, _ := goquery.NewDocumentFromReader(readerCurContents)
	listDom := contentsDom.Find(".articleList").Children()
	listLen := listDom.Length()

	for i := 1; i <= listLen; i++ {
		iStr := strconv.Itoa(i)
		page.Find(".articleList > li:nth-child(" + iStr + ") > a").Click()
		page.Back()
		time.Sleep(2 * time.Second)
		fmt.Println("kita")
	}
}
