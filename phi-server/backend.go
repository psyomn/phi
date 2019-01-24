/*
Copyright 2019 Simon Symeonidis (psyomn)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

const (
	dbName = "phidb.sqlite3"

	phiSchema = `CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username STRING UNIQUE,
  password STRING,
  salt BIGINT
)`
	insertUserSQL = `INSERT INTO users (username, password, salt) VALUES (?,?,?)`
	loginUserSQL  = `SELECT username, password, salt FROM users WHERE username = '?'`
)

var (
	dbHandle *sql.DB
	dbMutex  sync.Mutex
)

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

func encryptPasswordWithSalt(password string, salt int64) string {
	saltBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(saltBytes, uint64(salt))

	saltedPassword := append(saltBytes, password...)
	hashedPassword, err := bcrypt.GenerateFromPassword(saltedPassword, 14)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

func registerUser(username, password string, mutex *sync.Mutex) error {
	mutex.Lock()
	defer mutex.Unlock()

	hashedPassword, salt := encryptPassword(password)
	db, err := getDb()
	if err != nil {
		return fmt.Errorf("could not open db to store user: %v", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("could not start transaction: %v", err)
	}

	stmt, err := tx.Prepare(insertUserSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, hashedPassword, salt)
	if err != nil {
		tx.Rollback()

		if strings.Contains(err.Error(), "UNIQUE") {
			return errors.New("username has been taken")
		}

		return err
	}

	tx.Commit()
	return nil
}

func login(username, password string) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	db, err := getDb()
	if err != nil {
		return "", fmt.Errorf("could not open db to store user: %v", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return "", fmt.Errorf("could not start transaction: %v", err)
	}

	stmt, err := tx.Prepare(loginUserSQL)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	{ // Check if password matches
		var (
			dbUsername string
			dbPassword string
			dbSalt     int64
		)

		err = stmt.QueryRow("1").Scan(
			&dbUsername, &dbPassword, &dbSalt)

		if err != nil {
			return "", errors.New("user not found")
		}

		encryptedPass := encryptPasswordWithSalt(password, dbSalt)

		if encryptedPass != dbPassword {
			return "", errors.New("could not log in")
		}
	}

	userToken := make([]byte, 8)
	_, err = rand.Read(userToken)
	if err != nil {
		panic(err)
	}
	tokenInt := int64(binary.BigEndian.Uint64(userToken))
	tokenHex := fmt.Sprintf("%x", tokenInt)

	return tokenHex, nil
}
