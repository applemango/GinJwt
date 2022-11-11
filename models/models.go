package models

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenBlockList struct {
	Id        int    `json:"id"`
	UserId    int    `json:"userId"`
	Token     string `json:"token"`
	Signature string `json:"signature"`
}

func ConnectDB() error {
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		return err
	}

	create_user_query := "CREATE TABLE IF NOT EXISTS user ( id integer primary key autoincrement, username string unique, password string )"
	create_tokenblocklist_query := "CREATE TABLE IF NOT EXISTS tokenblocklist ( id integer primary key autoincrement, userid integer, token string, signature string )"

	db.Exec(create_user_query)
	_, err = db.Exec(create_tokenblocklist_query)

	if err != nil {
		fmt.Fprintf(os.Stderr, "err : createDB :%v\n", err)
		return err
	}
	DB = db
	return nil
}
