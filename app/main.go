package main

import (
	"database/sql"
	"strconv"

	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

var db sql.DB

type Donator struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	PostCode string `json:"postcode"`
}

func main() {
	fmt.Printf("Hello, World\n")

	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	sqlStmt := `
	create table IF NOT EXISTS foo (id integer not null primary key, name text);
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalln(err)
	}

	trans, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	_ = trans

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World\n")
	})

	e.GET("/users", func(c echo.Context) error {
		tx, err := db.Begin()
		defer tx.Commit()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}

		rows, err := tx.Query("SELECT * FROM foo;")

		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}

		donators := []Donator{}

		for rows.Next() {

			var name string
			var id int
			if rows.Scan(&id, &name) != nil {
				return c.JSON(http.StatusInternalServerError, nil)
			}

			donators = append(donators, Donator{
				Id:   id,
				Name: name,
			})
		}
		log.Println("Succesfully Retrieved", len(donators), " Donors")
		return c.JSON(http.StatusOK, donators)
	})

	e.POST("/users", func(c echo.Context) error {
		var donor Donator
		if err := c.Bind(&donor); err != nil {
			return err
		}
		RowsAffected, err := AddDonor(db, donor)
		if err != nil {
			return err
		}
		resp := struct {
			Ok           bool    `json:"ok"`
			Data         Donator `json:"data"`
			RowsAffected int     `json:"rowsAffected"`
		}{
			Data:         donor,
			Ok:           true,
			RowsAffected: RowsAffected,
		}
		return c.JSON(http.StatusOK, resp)
	})

	e.PUT("/users/:id", func(c echo.Context) error {
		var donor Donator
		if err := c.Bind(&donor); err != nil {
			return err
		}

		donor.Id, err = strconv.Atoi(c.Param("id"))

		RowsAffected, err := UpdateDonor(db, donor)
		if err != nil {
			return err
		}
		resp := struct {
			Ok           bool    `json:"ok"`
			Data         Donator `json:"data"`
			RowsAffected int     `json:"rowsAffected"`
		}{
			Data:         donor,
			Ok:           true,
			RowsAffected: RowsAffected,
		}
		return c.JSON(http.StatusOK, resp)
	})

	e.GET("/users/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Logger().Error(err)
			return err
		}

		donor, err := GetDonor(db, id)

		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusBadRequest, err)
		}

		log.Printf("Donor ID: %v\n", donor.Id)

		return c.JSON(200, donor)

	})

	e.Logger.Fatal(e.Start(":8080"))
}
