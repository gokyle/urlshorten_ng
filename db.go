package main

import (
	"bitbucket.org/taruti/pbkdf2"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

const dbFile = "data/urlshorten.db"

func dbConnect() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", dbFile)
	return
}

func lookupShortCode(sid string) (url string, err error) {
	db, err := dbConnect()
	if err != nil {
		return
	}
	defer db.Close()
	var rows *sql.Rows
	rows, err = db.Query("select url from shortened where sid=$1", sid)
	if err != nil {
		return
	}

	for rows.Next() {
		rows.Scan(&url)
	}
	err = rows.Err()
	return
}

func urlToSid(url string) (sid string, err error) {
	db, err := dbConnect()
	if err != nil {
		return
	}
	defer db.Close()
	var rows *sql.Rows
	rows, err = db.Query("select sid from shortened where url=?", url)
	if err != nil {
		return
	} else {
		for rows.Next() {
			rows.Scan(&sid)
		}
		err = rows.Err()
	}
	return
}

func getPassHash(username string) (ph pbkdf2.PasswordHash, err error) {
	db, err := dbConnect()
	if err != nil {
		return
	}
	defer db.Close()
	query := "select * from users where username=?"
	rows, err := db.Query(query, username)
	if err != nil {
		return
	}

	var user, hashed, salt string
	for rows.Next() {
		err = rows.Scan(&user, &hashed, &salt)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if err = rows.Err(); err != nil {
		return
	}
	ph = pbkdf2.PasswordHash{[]byte(salt), []byte(hashed)}
	return
}

func insertShortened(shorturl, url string) (err error) {
	db, err := dbConnect()
	if err != nil {
		return
	}
	defer db.Close()
	_, err = db.Exec("insert into shortened values (?, ?)",
		shorturl, url)
	return
}
