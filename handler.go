package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type handler struct {
	DB *sql.DB
}

func NewApp(newTable bool) *handler {
	h := conToDB()

	if newTable {
		h.createTable()
	}
	return h
}

func conToDB() *handler {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	return &handler{db}

}
func (h *handler) createTable() {
	sqlcreateTb := `CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);`

	_, err := h.DB.Exec(sqlcreateTb)
	if err != nil {
		log.Fatal("can't create table", err)
	}
	fmt.Println("create table success")
}
