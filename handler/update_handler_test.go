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

func TestUpdateExpenseByIDUnit(t *testing.T) {

	//Arrange
	reqBody := `
	{"id": 1,
    "title": "apple smoothie",
    "amount": 89,
    "note": "no discount",
    "tags": ["beverage"]}`

	req := httptest.NewRequest(http.MethodPut, "/expenses/1", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	db, mock, err := sqlmock.New()
	newsMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(1, "apple smoothie", 89, "no discount", `{"beverage"}`)
	mock.ExpectPrepare("UPDATE expenses SET title=\\$2, amount=\\$3, note=\\$4, tags=\\$5 WHERE id=\\$1 RETURNING id, title, amount, note, tags").ExpectQuery().WithArgs(1, "apple smoothie", 89.0, "no discount", `{"beverage"}`).WillReturnRows(newsMockRows)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	h := handler{db}
	c := echo.New().NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	expected := "{\"id\":1,\"title\":\"apple smoothie\",\"amount\":89,\"note\":\"no discount\",\"tags\":[\"beverage\"]}"
	//Act
	err = h.UpdateExpenseHandler(c)

	//Assert
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}

}
