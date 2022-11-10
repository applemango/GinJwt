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

func ConnectDB() error {
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		return err
	}

	query := "CREATE TABLE IF NOT EXISTS user ( id integer primary key autoincrement, username string unique, password string )"

	_, err = db.Exec(query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err : createDB :%v\n", err)
		return err
	}
	DB = db
	return nil
}

func GetUserFromId(id int) (User, error) {
	stmt, err := DB.Prepare("SELECT id, username, password from user WHERE id = ?")
	if err != nil {
		fmt.Fprintf(os.Stderr, "err : getDB :%v\n", err)
		return User{}, err
	}
	user := User{}
	sqlErr := stmt.QueryRow(id).Scan(&user.Id, &user.Username, &user.Password)
	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return User{}, nil
		}
		return User{}, sqlErr
	}
	return user, nil
}

func GetUserFromUsername(username string) (User, error) {
	row := DB.QueryRow(`SELECT * from user WHERE username = ?`, username)
	user := User{}
	var err error
	if err = row.Scan(&user.Id, &user.Username, &user.Password); err == sql.ErrNoRows {
		fmt.Println("not found")
		return User{}, err
	}
	fmt.Println(user.Id, user.Username, user.Password)
	return user, nil
	/*user := User{}
	sqlErr := stmt.QueryRow(username).Scan(&user.Id, &user.Username, &user.Password)
	fmt.Println(sqlErr)
	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return User{}, nil
		}
		return User{}, sqlErr
	}*/
	/*for stmt.Next() {
		sqlErr := stmt.Scan(&user.Id, &user.Username, &user.Password)
		if sqlErr != nil {
			return User{}, sqlErr
		}
	}
	fmt.Println(user, err)*/
	//return user, nil
}

func GetUserFromLastInsert() (User, error) {
	stmt, err := DB.Query("SELECT id, username, password from user WHERE id = last_insert_rowid()")

	if err != nil {
		fmt.Fprintf(os.Stderr, "err : getDB :%v\n", err)
		return User{}, err
	}

	defer stmt.Close()
	user := User{}

	for stmt.Next() {
		sqlErr := stmt.Scan(&user.Id, &user.Username, &user.Password)
		if sqlErr != nil {
			return User{}, sqlErr
		}
	}
	return user, nil
}

func InsertUser(newUser User) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO user ( username, password ) VALUES ( ?, ? )")
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: addUrl :%v\n", err)
		return false, err
	}

	_, err = stmt.Exec(newUser.Username, newUser.Password)
	if err != nil {
		return false, err
	}

	tx.Commit()
	return true, nil
}