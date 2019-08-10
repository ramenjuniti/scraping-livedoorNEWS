# scraping-livedoorNEWS

## About

[livedoor NEWS](http://news.livedoor.com/)から以下の 3 つの情報を取得し、csv ファイルを生成する

- 記事タイトル
- 記事本文
- 記事の 3 文要約（ざっくり言うと）

**記事がある限り取得し続けます**

## Dependency

- goquery: [github](https://github.com/PuerkitoBio/goquery), [GoDoc](https://godoc.org/github.com/PuerkitoBio/goquery)
- agouti: [github](https://github.com/sclevine/agouti), [GoDoc](https://godoc.org/github.com/sclevine/agouti), [公式](https://agouti.org/)

## Usage

```
brew cask install chromedriver
git clone git@github.com:ramenjuniti/scraping-livedoorNEWS.git
cd scraping-livedoorNEWS
go run main.go sample.csv
```

## License

This software is released under the MIT License, see LICENSE.
