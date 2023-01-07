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
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
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

	_, err = db.Exec(createTb)
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
func getIdExpenseHandlers(c echo.Context) error {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, model.Errmsg{Message: "Id is not found value"})
	}
	result := model.Resbody{}
	row := db.QueryRow("SELECT id, title, amount, note, tags FROM expenses WHERE id=$1", id)

	err1 := row.Scan(&result.Id, &result.Title, &result.Amount, &result.Note, pq.Array(&result.Tags))
	if err1 != nil {
		return c.JSON(http.StatusInternalServerError, model.Errmsg{Message: err.Error()})
	}

	resBody := model.Resbody{
		Id:     result.Id,
		Title:  result.Title,
		Amount: result.Amount,
		Note:   result.Note,
		Tags:   result.Tags,
	}
	defer db.Close()

	return c.JSON(http.StatusOK, resBody)

}
func getAllexpensesHandler(c echo.Context) error {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, model.Errmsg{Message: "Id is not found value"})
	}
	result := model.Resbody{}
	row := db.QueryRow("SELECT id, title, amount, note, tags FROM expenses WHERE id=$1", id)

	err1 := row.Scan(&result.Id, &result.Title, &result.Amount, &result.Note, pq.Array(&result.Tags))
	if err1 != nil {
		return c.JSON(http.StatusInternalServerError, model.Errmsg{Message: err.Error()})
	}

	resBody := model.Resbody{
		Id:     result.Id,
		Title:  result.Title,
		Amount: result.Amount,
		Note:   result.Note,
		Tags:   result.Tags,
	}
	defer db.Close()

	return c.JSON(http.StatusOK, resBody)

}
func updateExpenseHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, model.Errmsg{Message: "Id is not found value"})
	}
	var body model.Reqbody
	err := c.Bind(&body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Errmsg{Message: err.Error()})
	}
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	result := model.Resbody{}
	row := db.QueryRow("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id=$1 RETURNING *",
		id, body.Title, body.Amount, body.Note, pq.Array(body.Tags))

	err = row.Scan(&result.Id, &result.Title, &result.Amount, &result.Note, pq.Array(&result.Tags))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Errmsg{Message: err.Error()})

	}
	resBody := model.Resbody{
		Id:     result.Id,
		Title:  result.Title,
		Amount: result.Amount,
		Note:   result.Note,
		Tags:   result.Tags,
	}
	return c.JSON(http.StatusOK, resBody)

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
	e.GET("/expenses/:id", getIdExpenseHandlers)
	e.GET("/expenses", getAllexpensesHandler)
	e.PUT("/expenses/:id", updateExpenseHandler)
	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
