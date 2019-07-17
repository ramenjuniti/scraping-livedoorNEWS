package main

import (
	"encoding/csv"
	"errors"
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

		fmt.Println(listLen)
		fmt.Println(dataCount)
		// 記事のリストをぐるぐる
		for i := 1; i <= listLen; i++ {
			iStr := strconv.Itoa(i)
			if _, err := page.Find(".articleList > li:nth-child(" + iStr + ") > a:visited").Text(); err == nil {
				continue
			}

			// 記事のタイトルと要約へ
			page.Find(".articleList > li:nth-child(" + iStr + ") > a").Click()
			time.Sleep(sleepTime * time.Second)

			// タイトル取得
			articleTitle, err := page.Find(".topicsTtl > a").Text()
			if err != nil {
				page.Back()
				time.Sleep(sleepTime * time.Second)
				continue
			}

			// 要約取得
			summary, err := page.FindByClass("summaryList").Text()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				page.Back()
				time.Sleep(sleepTime * time.Second)
				continue
			}

			// 要約が3文かどうかチェック
			summaryList := strings.Split(summary, "\n")
			articleMoreButton := page.Find(".articleMore > a")
			_, err = articleMoreButton.Text()
			if len(summaryList) != 3 || err != nil {
				if len(summaryList) != 3 {
					fmt.Fprintln(os.Stderr, errors.New("there are less than 3 summaries"))
				}
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				page.Back()
				time.Sleep(sleepTime * time.Second)
				continue
			}

			// 記事の本文へ
			articleMoreButton.Click()
			time.Sleep(sleepTime * time.Second)

			// 記事の本文取得
			articleBody, err := page.Find(".articleBody > span").Text()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				page.Back()
				time.Sleep(sleepTime * time.Second)
				page.Back()
				time.Sleep(sleepTime * time.Second)
				continue
			}

			// csvに書き込み
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

			// 本文 -> タイトル＆要約
			page.Back()
			time.Sleep(sleepTime * time.Second)

			// タイトル＆要約 -> 記事リスト
			page.Back()
			time.Sleep(sleepTime * time.Second)
		}

		// 次の記事リストへ
		nextPage := page.Find(".next > a")
		_, err := nextPage.Text()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}

		nextPage.Click()
		time.Sleep(sleepTime * time.Second)
	}
}
