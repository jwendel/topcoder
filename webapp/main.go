package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type A struct {
	Name string
	Id   int
}

func main() {
	listenAddr := flag.String("listen", ":8080", "Host/port to listen on")
	flag.Parse()

	fmt.Println(listenAddr)
	// http.ListenAndServe(listenAddr, handler)
	db, err := sqlx.Open("postgres", "user=gotest password=gotest dbname=testdb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Queryx("SELECT * FROM testtable")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var ss A
		rows.StructScan(&ss)
		log.Println(ss)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
