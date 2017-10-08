package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func LoadDB(user, pass, dbName string) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, pass, dbName))
	if err != nil {
		return nil, fmt.Errorf("loaddb failure: %v", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("loaddb ping failure: %v", err)
	}
	return db, nil
}
