//go:build unit
// +build unit

package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetAllExpenseUnit(t *testing.T) {

	//Arrange
	req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	newsMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(1, "strawberry smoothie", "79", "night market promotion discount 10 bath", `{"food","beverage"}`)
	db, mock, err := sqlmock.New()
	mock.ExpectQuery("SELECT (.+) FROM expenses").WillReturnRows(newsMockRows)
	// fmt.Println(newsMockRows)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	h := handler{db}
	c := echo.New().NewContext(req, rec)
	expected := "[{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}]"

	//Act
	err = h.GetExpensesHandler(c)
	//Arrange
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}

}
