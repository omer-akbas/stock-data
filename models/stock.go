package models

//Son işlem fiyatı = last price
//Alış fiyatı = bid
//Satış fiyatı = ask
//Önceki kapanış fiyatı = PreviousPrice
type Stock struct {
	LastPrice, PreviousPrice, Bid, Ask float64
	Name, Code                         string
}
