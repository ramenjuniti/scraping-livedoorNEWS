package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sclevine/agouti"
)

func main() {
	flag.Parse()
	arg := flag.Arg(0)

	file, err := os.OpenFile(arg, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("Failed to open file: %v", err)
	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		log.Printf("Failed to clear file: %v", err)
	}

	writer := csv.NewWriter(file)
	writer.Write([]string{"", "title", "body", "summary1", "summary2", "summary3"})
	writer.Flush()

	url := "http://news.livedoor.com/topics/category/main/"
	driver := agouti.ChromeDriver()

	err = driver.Start()
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

	for {
		listDom := contentsDom.Find(".articleList").Children()
		listLen := listDom.Length()

		for i := 1; i <= listLen; i++ {
			fmt.Printf("%v 番目の記事の情報を取得します\n", i)
			iStr := strconv.Itoa(i)
			page.Find(".articleList > li:nth-child(" + iStr + ") > a").Click()
			time.Sleep(5 * time.Second)

			summary, err := page.FindByClass("summaryList").Text()
			if err == nil {
				summaryList := strings.Split(summary, "\n")
				articleMoreButton := page.Find(".articleMore > a")
				_, err := articleMoreButton.Text()
				if len(summaryList) == 3 && err == nil {
					articleMoreButton.Click()
					time.Sleep(5 * time.Second)

					articleTitle, err := page.Find(".articleTtl").Text()
					articleBody, err := page.Find(".articleBody > span").Text()

					if err == nil {
						writer.Write([]string{
							articleTitle,
							articleBody,
							summaryList[0],
							summaryList[1],
							summaryList[2],
						})
						writer.Flush()
					}
					page.Back()
				}
			}
			page.Back()
			time.Sleep(5 * time.Second)
		}
		nextPage := page.Find(".next > a")
		_, err := nextPage.Text()
		if err != nil {
			break
		}
		nextPage.Click()
		time.Sleep(5 * time.Second)
	}
}
