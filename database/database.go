package database

import (
	"database/sql"
	"log"
	_ "modernc.org/sqlite" // import drivers only
)

var DB *sql.DB

func InitDatabase() {
	var err error
	DB, err = sql.Open("sqlite", "./blog.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	createTable()
}

func createTable() {
	rawCreateStatement := `CREATE TABLE IF NOT EXISTS blog (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"author" TEXT,
		"title" TEXT,
		"text" TEXT,
		"creation_time" TEXT
		);`

	createStatement, err := DB.Prepare(rawCreateStatement)
	if err != nil {
		log.Fatal("Failed to create database table:", err)
	}

	_, execErr := createStatement.Exec()
	if execErr != nil {
		log.Fatal("Failed to execute create table statement:", execErr)
	}
}

func CloseDatabase() {
	DB.Close()
}
