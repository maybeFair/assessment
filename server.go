package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/maybeFair/assessment/model"
)

func init() {
	sqldb, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	createTb := `CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);`

	_, err = sqldb.Exec(createTb)
	if err != nil {
		log.Fatal("can't create table", err)
	}
	fmt.Println("create table success")

}
func createExpenseHandler(c echo.Context) error {

	var tbId int
	var body model.Reqbody
	err := c.Bind(&body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Errmsg{Message: err.Error()})
	}
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	row := db.QueryRow("INSERT INTO expenses (id, title,amount,note,tags)  values (default, $1, $2, $3, $4) RETURNING id", body.Title, body.Amount, body.Note, pq.Array(body.Tags))
	err = row.Scan(&tbId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Errmsg{Message: err.Error()})
	}
	resbody := model.Resbody{
		Id:     tbId,
		Title:  body.Title,
		Amount: body.Amount,
		Note:   body.Note,
		Tags:   body.Tags}

	return c.JSON(http.StatusCreated, resbody)

}
func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "November",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == "10, 2009", nil
		},
	}))
	e.POST("/expenses", createExpenseHandler)
	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
