package models

//Son işlem fiyatı = last price
//Alış fiyatı = bid
//Satış fiyatı = ask
//Önceki kapanış fiyatı = PreviousPrice
type Stock struct {
	LastPrice, PreviousPrice, Bid, Ask float64
	Name, Code                         string
}

func (s *Stock) Insert() error {
	db := dbConnect()
	defer db.Close()

	var (
		id, count int
		statement string
	)

	_ = db.QueryRow("SELECT COUNT(*) FROM stock WHERE code = ?", s.Code).Scan(&count)
	if count == 0 { //first data
		lastInsert, err := db.Exec("INSERT INTO stock(name, code) values(?, ?)", s.Name, s.Code)
		if err != nil {
			return err
		}

		returnId, err := lastInsert.LastInsertId()
		if err != nil {
			return err
		}
		id = int(returnId)
		statement = "-"
	} else { //not the first
		var lstBid float64

		err := db.QueryRow("SELECT id FROM stock WHERE code = ? LIMIT 1", s.Code).Scan(&id)
		if err != nil {
			return err
		}

		err = db.QueryRow("SELECT bid FROM price WHERE  stockId = ? LIMIT 1 ORDER BY id DESC", id).Scan(&lstBid)
		if err != nil {
			return err
		}
		if lstBid > s.Bid {
			statement = "high"
		} else if lstBid < s.Bid {
			statement = "low"
		} else {
			statement = "-"
		}
	}

	_, err := db.Exec("INSERT INTO price(lastPrice, previousPrice, bid, ask, status, stockId) values(?, ?, ?, ?, ?, ?)", s.LastPrice, s.PreviousPrice, s.Bid, s.Ask, statement, id)
	if err != nil {
		return err
	}

	return nil
}
