package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
)

var reqbodyTest = `{
	"title": "buy a new phone",
	"amount": 39000,
	"note": "buy a new phone",
	"tags": ["gadget", "shopping"]
}`

var resbodyTest = `{
	"id": 1,
	"title": "buy a new phone",
	"amount": 39000,
	"note": "buy a new phone",
	"tags": ["gadget", "shopping"]
}`

func TestcreateExpenseHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(reqbodyTest))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal("cannot mock db", err.Error())
	}
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mockSqlCmd := "INSERT INTO expenses (id, title, amount, note, tags) values (default, $1, $2, $3, $4) RETURNING id"

	mock.ExpectQuery(regexp.QuoteMeta(mockSqlCmd)).WithArgs("strawberry smoothie",
		float64(79),
		"night market promotion discount 10 bath",
		pq.Array([]string{"food", "beverage"})).WillReturnRows(rows)

	c := e.NewContext(req, rec)

	if assert.NoError(t, createExpenseHandler(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, resbodyTest, strings.TrimSpace(rec.Body.String()))
	}

}
