package main

import (
	"database/sql"
	"fmt"
)

func UpdateDonor(db *sql.DB, donor Donator) (int, error) {
	tx, err := db.Begin()
	defer tx.Commit()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.Exec("UPDATE foo SET name = ? WHERE id = ?", donor.Name, donor.Id)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := stmt.RowsAffected()

	return int(rowsAffected), nil
}

func GetDonor(db *sql.DB, id int) (*Donator, error) {
	tx, err := db.Begin()
	defer tx.Commit()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Query("SELECT * FROM foo WHERE id = ?", id)

	if err != nil {
		return nil, err
	}

	if !stmt.Next() {
		return nil, fmt.Errorf("Failed To Find Donor With ID %v", id)
	}
	var name string
	stmt.Scan(&id, &name)

	return &Donator{
		Name: name,
		Id:   id,
	}, nil
}

func AddDonor(db *sql.DB, d Donator) (int, error) {

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	stmt, err := db.Exec("INSERT INTO foo(name) values(?)", d.Name)
	affectedCount, err := stmt.RowsAffected()

	if err != nil {
		return 0, err
	}

	defer tx.Commit()
	return int(affectedCount), nil
}
