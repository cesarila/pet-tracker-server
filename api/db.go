package main

import (
	"context"
	"database/sql"
	"strings"

	_ "modernc.org/sqlite"
)

type Pet struct {
	Name   string
	Inside bool
}

type PetDbRow struct {
	id int
	Pet
}

var sqliteDB *sql.DB

func initDatabase(dbPath string) (*sql.DB, error) {
	var err error
	sqliteDB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	_, err = sqliteDB.ExecContext(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS t_pets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			inside INTEGER NOT NULL DEFAULT 1
		)`,
	)
	if err != nil {
		return nil, err
	}
	return sqliteDB, nil
}

func getPets() ([]PetDbRow, error) {
	var pets []PetDbRow

	rows, err := sqliteDB.QueryContext(
		context.Background(),
		`SELECT * FROM t_pets`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var pet PetDbRow
		if err := rows.Scan(
			&pet.id, &pet.Name, &pet.Inside,
		); err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}
	return pets, err
}

func addPet(p *Pet) (int64, error) {
	result, err := sqliteDB.ExecContext(
		context.Background(),
		`INSERT INTO t_pets (name, inside) VALUES (?,?);`, p.Name, p.Inside,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: t_pets.name") == true {
			return 0, nil
		} else {
			return -1, err
		}
	}
	rowsAffected, err := result.RowsAffected()
	// id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return rowsAffected, nil
}

func updatePetStatus(p *Pet) (int64, error) {
	result, err := sqliteDB.ExecContext(
		context.Background(),
		`UPDATE t_pets
			SET inside = (?)
			WHERE name = (?);`, p.Inside, p.Name,
	)
	rowCount, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}
	return rowCount, err
}

func deletePet(p *Pet) (int64, error) {
	result, err := sqliteDB.ExecContext(
		context.Background(),
		`DELETE FROM t_pets
			WHERE name = (?);`, p.Name,
	)
	rowCount, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}
	return rowCount, err
}
