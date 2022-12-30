package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eskimoburger/assessment/handler"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func AuthorizationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader != os.Getenv("Authorization") {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		}
		return next(c)
	}
}

func main() {
	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("can't create table", err)
	}
	createTb := `
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);
	`
	_, err = db.Exec(createTb)
	if err != nil {
		log.Fatal("can't create table", err)
	}
	h := handler.NewApplication(db)
	e := echo.New()
	e.Use(AuthorizationMiddleware)
	e.POST("/expenses", h.CreateExpenseHandler)
	e.GET("/expenses/:id", h.GetExpenseByIDHandler)
	e.PUT("/expenses/:id", h.UpdateExpenseHandler)
	e.GET("/expenses", h.GetExpensesHandler)

	log.Fatal(e.Start(os.Getenv("PORT")))

}
