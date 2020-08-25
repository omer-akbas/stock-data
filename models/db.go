package models

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func dbConnect() *sql.DB {
	var (
		host     string
		port     string
		database string
		user     string
		password string
	)

	flag.StringVar(&host, "host", "93.187.203.193", "host (ip or name)")
	flag.StringVar(&port, "port", "3306", "port")
	flag.StringVar(&database, "database", "mynet", "database name")
	flag.StringVar(&user, "user", "mynet", "username")
	flag.StringVar(&password, "password", "#871mtkV", "password")

	flag.Parse()
	connQuery := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)
	db, err := sql.Open("mysql", connQuery)
	if err != nil {
		log.Println("db connection err: ", err.Error())
		return nil
	}
	return db
}
