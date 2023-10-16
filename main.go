package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "jamabot"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		rows, err := db.QueryContext(c.Request().Context(), "select u.id,u.name,u.email from users u")
		if err != nil {
			return err
		}
		defer rows.Close()

		type row = struct {
			Id    *int32  `json:"id,omitempty"`
			Name  *string `json:"name,omitempty"`
			Email *string `json:"email,omitempty"`
		}

		var res = []row{}
		for rows.Next() {
			var row = row{}
			if err := rows.Scan(&row.Id, &row.Name, &row.Email); err != nil {
				c.Logger().Error(err)
				return err
			}
			res = append(res, row)
		}
		if err = rows.Err(); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, res)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
