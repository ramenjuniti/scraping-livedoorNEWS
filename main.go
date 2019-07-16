package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sclevine/agouti"
)

const sleepTime = 5
const url = "https://news.livedoor.com/topics/category/dom/"

func replace(t string) string {
	r := strings.NewReplacer("\n", "", ",", "、")
	return r.Replace(t)
}

func main() {
	flag.Parse()
	arg := flag.Arg(0)

	file, err := os.OpenFile(arg, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to clear file: %v", err)
		os.Exit(1)
	}

	writer := csv.NewWriter(file)
	writer.Write([]string{"", "title", "body", "summary1", "summary2", "summary3"})
	writer.Flush()

	driver := agouti.ChromeDriver()

	err = driver.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start driver: %v", err)
		os.Exit(1)
	}
	defer driver.Stop()

	page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open page: %v", err)
		os.Exit(1)
	}

	err = page.Navigate(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to navigate: %v", err)
		os.Exit(1)
	}

	curContentsDom, err := page.HTML()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get html: %v", err)
		os.Exit(1)
	}

	readerCurContents := strings.NewReader(curContentsDom)
	contentsDom, _ := goquery.NewDocumentFromReader(readerCurContents)

	var dataCount int

	for {
		listDom := contentsDom.Find(".articleList").Children()
		listLen := listDom.Length()

		for i := 1; i <= listLen; i++ {
			iStr := strconv.Itoa(i)
			if _, err := page.Find(".articleList > li:nth-child(" + iStr + ") > a:visited").Text(); err == nil {
				continue
			}
			page.Find(".articleList > li:nth-child(" + iStr + ") > a").Click()
			time.Sleep(sleepTime * time.Second)
			summary, err := page.FindByClass("summaryList").Text()

			if err == nil {
				summaryList := strings.Split(summary, "\n")
				articleMoreButton := page.Find(".articleMore > a")
				_, err := articleMoreButton.Text()

				if len(summaryList) == 3 && err == nil {
					articleMoreButton.Click()
					time.Sleep(sleepTime * time.Second)
					articleTitle, err := page.Find(".articleTtl").Text()
					articleBody, err := page.Find(".articleBody > span").Text()

					if err == nil {
						writer.Write([]string{
							replace(articleTitle),
							replace(articleBody),
							replace(summaryList[0]),
							replace(summaryList[1]),
							replace(summaryList[2]),
						})
						writer.Flush()
						dataCount++
						fmt.Printf("現在 %v 個の記事を取得済みです\n", dataCount)
					}
					page.Back()
				}
			}
			page.Back()
			time.Sleep(sleepTime * time.Second)
		}
		nextPage := page.Find(".next > a")
		_, err := nextPage.Text()
		if err != nil {
			break
		}
		nextPage.Click()
		time.Sleep(sleepTime * time.Second)
	}
}
