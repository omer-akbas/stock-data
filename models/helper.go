package models

import (
	"log"
	"strconv"
	"strings"
	"time"
)

func toFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 8)
	return f
}

func chronometer(startTime time.Time) {
	endTime := time.Since(startTime)
	log.Println("startTime: ", startTime, "endTime: ", endTime, "=========> ", shortDuration(endTime))
}

func shortDuration(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}
