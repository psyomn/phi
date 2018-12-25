package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

const (
	phiSchema = `
CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username STRING UNIQUE,
  password STRING,
  salt BIGINT
)
`
	insertUserSQL = `
INSERT INTO users (username, password, salt) VALUES (?,?,?)
`
)

const dbName = "phidb.sqlite3"

func init() {
	createDbIfNotExist()
}

func getDb() (*sql.DB, error) {
	return sql.Open("sqlite3", dbName)
}

func createDbIfNotExist() {
	if _, err := os.Stat(dbName); !os.IsNotExist(err) {
		return
	}

	db, err := getDb()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err != nil {
		panic(fmt.Sprintf("could not open db: %v", err))
	}

	_, err = db.Exec(phiSchema)
	if err != nil {
		panic(fmt.Sprintf("could not create schema: %v", err))
	}
}

func encryptPassword(password string) (string, int64) {
	saltBytes := make([]byte, 8)
	_, err := rand.Read(saltBytes)
	if err != nil {
		panic(err)
	}
	salt := int64(binary.BigEndian.Uint64(saltBytes))

	saltedPassword := append(saltBytes, password...)
	hashedPassword, err := bcrypt.GenerateFromPassword(saltedPassword, 14)

	if err != nil {
		panic(err)
	}

	return string(hashedPassword), salt
}

func registerUser(username, password string) {
	hashedPassword, salt := encryptPassword(password)
	db, err := getDb()

	if err != nil {
		panic(fmt.Sprintf("could not open db to store user: %v", err))
	}

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("could not start transaction: %v", err))
	}

	stmt, err := tx.Prepare(insertUserSQL)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, hashedPassword, salt)
	if err != nil {
		panic(err)
	}
	tx.Commit()
}
