package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/maybeFair/assessment/model"
)

func createExpenseHandler(c echo.Context) error {
	var body model.Reqbody
	err := c.Bind(&body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Errmsg{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, body)

}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("/expenses", createExpenseHandler)
	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))
}
