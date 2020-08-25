package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/omer-akbas/stock-data/utils"
	"github.com/robfig/cron/v3"
)

func main() {
	scrapperStart()
	c := cron.New()
	c.AddFunc("@every 1h0m0s", scrapperStart)
	c.Start()
	time.Sleep(240 * time.Hour)
	c.Stop()
}

func scrapperStart() {

	defer utils.Chronometer(time.Now())

	urls := urlList()

	stocks := []Stock{}

	var wg sync.WaitGroup
	var lock sync.Mutex

	for i, url := range urls {
		if i%10 == 0 {
			time.Sleep(3 * time.Second)
		}
		wg.Add(1)
		go urlVisit(url, &stocks, &wg, &lock)
	}
	wg.Wait()

	fmt.Println("--------------------------------------")
	// fmt.Println(stocks)
	fmt.Println("ADET: ", len(stocks))
}

//Son işlem fiyatı = last price
//Alış fiyatı = bid
//Satış fiyatı = ask
//Önceki kapanış fiyatı = PreviousPrice
type Stock struct {
	LastPrice, PreviousPrice, Bid, Ask float64
	Name, Code                         string
}

func urlVisit(url string, stocks *[]Stock, wg *sync.WaitGroup, lock *sync.Mutex) {
	defer wg.Done()
	c := colly.NewCollector()

	var stock Stock

	c.OnHTML("body", func(e *colly.HTMLElement) {
		stock = Stock{
			Name:          e.ChildText("section > div.row > div.col-12.col-md-8.col-content > div:nth-child(3) > div > div:nth-child(1) > div.data-detay-page-heading > div:nth-child(1) > div.col-9.flex.align-items-center > h1"),
			Code:          e.ChildText("section > div.row > div.col-12.col-md-8.col-content > div:nth-child(3) > div > div:nth-child(1) > div.data-detay-page-heading > div:nth-child(1) > div.col-9.flex.align-items-center > span"),
			LastPrice:     utils.ToFloat(e.ChildText("section > div.row > div.col-12.col-md-8.col-content > div:nth-child(3) > div > div:nth-child(1) > div.p-3 > div.flex-list-2-col.flex.justify-content-between > ul:nth-child(1) > li:nth-child(1) > span:nth-child(2)")),
			PreviousPrice: utils.ToFloat(e.ChildText("section > div.row > div.col-12.col-md-8.col-content > div:nth-child(3) > div > div:nth-child(1) > div.p-3 > div.flex-list-2-col.flex.justify-content-between > ul:nth-child(2) > li:nth-child(1) > span:nth-child(2)")),
			Bid:           utils.ToFloat(e.ChildText("section > div.row > div.col-12.col-md-8.col-content > div:nth-child(3) > div > div:nth-child(1) > div.p-3 > div.flex-list-2-col.flex.justify-content-between > ul:nth-child(1) > li:nth-child(2) > span:nth-child(2)")),
			Ask:           utils.ToFloat(e.ChildText("section > div.row > div.col-12.col-md-8.col-content > div:nth-child(3) > div > div:nth-child(1) > div.p-3 > div.flex-list-2-col.flex.justify-content-between > ul:nth-child(1) > li:nth-child(3) > span:nth-child(2)")),
		}
	})

	c.Visit(url)

	lock.Lock()
	*stocks = append(*stocks, stock)
	lock.Unlock()

	fmt.Println(stock)
}

func urlList() []string {
	c := colly.NewCollector()

	urlList := []string{}

	c.OnHTML("body > section > div.row > div.col-12.col-md-8.col-content > div:nth-child(3) > div > div > table > tbody > tr", func(e *colly.HTMLElement) {
		urlList = append(urlList, e.Request.AbsoluteURL(e.ChildAttr("a", "href")))
	})

	c.Visit("https://finans.mynet.com/borsa/hisseler")

	return urlList
}
