package main

import (
	"time"

	"github.com/omer-akbas/stock-data/models"
	"github.com/robfig/cron/v3"
)

func main() {
	target := models.Target{Url: "https://finans.mynet.com/borsa/hisseler"}
	target.ScrapperStart()
	c := cron.New()
	c.AddFunc("@every 1h0m0s", target.ScrapperStart)
	c.Start()
	time.Sleep(240 * time.Hour)
	c.Stop()
}
