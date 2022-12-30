package handler

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) GetExpensesHandler(c echo.Context) error {

	rows, err := h.DB.Query("SELECT * FROM expenses")
	if err != nil {
		log.Fatal("can't prepare query all expense statement", err)
	}
	var exs = []Expense{}
	for rows.Next() {
		ex := Expense{}
		err = rows.Scan(&ex.ID, &ex.Title, &ex.Amount, &ex.Note, pq.Array(&ex.Tags))
		if err != nil {
			log.Fatal("Can't Scan row into variable", err)
		}
		exs = append(exs, ex)
	}
	return c.JSON(http.StatusOK, exs)

}
