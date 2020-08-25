package models

import (
	"fmt"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

type Target struct {
	Url string
}

func (t *Target) ScrapperStart() {

	defer chronometer(time.Now())

	urls := t.urlList()

	stocks := []Stock{}

	var wg sync.WaitGroup
	var lock sync.Mutex

	for i, url := range urls {
		if i%10 == 0 {
			time.Sleep(3 * time.Second)
		}
		wg.Add(1)
		go t.urlVisit(url, &stocks, &wg, &lock)
	}
	wg.Wait()

	fmt.Println("--------------------------------------")
	// fmt.Println(stocks)
	fmt.Println("ADET: ", len(stocks))
}

func (t *Target) urlList() []string {
	c := colly.NewCollector()

	urlList := []string{}

	c.OnHTML("body > section > div.row > div.col-12.col-md-8.col-content > div:nth-child(3) > div > div > table > tbody > tr", func(e *colly.HTMLElement) {
		urlList = append(urlList, e.Request.AbsoluteURL(e.ChildAttr("a", "href")))
	})

	// c.Visit("https://finans.mynet.com/borsa/hisseler")
	c.Visit(t.Url)

	return urlList
}

func (t *Target) urlVisit(url string, stocks *[]Stock, wg *sync.WaitGroup, lock *sync.Mutex) {
	defer wg.Done()
	c := colly.NewCollector()

	var stock Stock

	c.OnHTML("body", func(e *colly.HTMLElement) {
		stock = Stock{
			Name:          e.ChildText("section > div.row > div.col-12.col-md-8.col-content > div:nth-child(3) > div > div:nth-child(1) > div.data-detay-page-heading > div:nth-child(1) > div.col-9.flex.align-items-center > h1"),
			Code:          e.ChildText("section > div.row > div.col-12.col-md-8.col-content > div:nth-child(3) > div > div:nth-child(1) > div.data-detay-page-heading > div:nth-child(1) > div.col-9.flex.align-items-center > span"),
			LastPrice:     toFloat(e.ChildText("section > div.row > div.col-12.col-md-8.col-content > div:nth-child(3) > div > div:nth-child(1) > div.p-3 > div.flex-list-2-col.flex.justify-content-between > ul:nth-child(1) > li:nth-child(1) > span:nth-child(2)")),
			PreviousPrice: toFloat(e.ChildText("section > div.row > div.col-12.col-md-8.col-content > div:nth-child(3) > div > div:nth-child(1) > div.p-3 > div.flex-list-2-col.flex.justify-content-between > ul:nth-child(2) > li:nth-child(1) > span:nth-child(2)")),
			Bid:           toFloat(e.ChildText("section > div.row > div.col-12.col-md-8.col-content > div:nth-child(3) > div > div:nth-child(1) > div.p-3 > div.flex-list-2-col.flex.justify-content-between > ul:nth-child(1) > li:nth-child(2) > span:nth-child(2)")),
			Ask:           toFloat(e.ChildText("section > div.row > div.col-12.col-md-8.col-content > div:nth-child(3) > div > div:nth-child(1) > div.p-3 > div.flex-list-2-col.flex.justify-content-between > ul:nth-child(1) > li:nth-child(3) > span:nth-child(2)")),
		}
	})

	c.Visit(url)

	lock.Lock()
	*stocks = append(*stocks, stock)
	lock.Unlock()

	fmt.Println(stock)
}
