package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) UpdateExpenseHandler(c echo.Context) error {
	string_id := c.Param("id")

	id, err := strconv.Atoi(string_id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "invalid ID: " + err.Error()})
	}

	ex := Expense{}
	err = c.Bind(&ex)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	stmt, err := h.DB.Prepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id=$1 RETURNING id, title, amount, note, tags")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Can't prepare query expense statement:" + err.Error()})
	}

	updated_ex := Expense{}

	row := stmt.QueryRow(id, ex.Title, ex.Amount, ex.Note, pq.Array(ex.Tags))
	err = row.Scan(&updated_ex.ID, &updated_ex.Title, &updated_ex.Amount, &updated_ex.Note, pq.Array(&updated_ex.Tags))

	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "Expense not found"})
	case nil:
		return c.JSON(http.StatusOK, updated_ex)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "Can't scan expense:" + err.Error() + fmt.Sprintf("Error scanning value: %v", updated_ex)})
	}

}
