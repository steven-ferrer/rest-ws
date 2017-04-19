package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func main() {
	db, err := sql.Open("sqlite3", "../users.db")
	DB = db
	if err != nil {
		log.Fatal("Problem opening database file: ", err.Error())
	}

	getUsers()

	for i := 1; i < 6; i++ {
		deleteUser(i)
	}

	//close the database
	DB.Close()
}

func createUser(username, departname string, created time.Time) {
	//insert
	stmt, err := DB.Prepare("INSERT INTO userinfo(username, departname, created) VALUES(?,?,?)")

	fmt.Println(stmt)
	res, err := stmt.Exec("steven", "development", "2017-04-29")
	if err != nil {
		log.Fatal(err.Error())
	}

	//get the id of the last inserted row
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Last Inserted ID: ", id)
	//update
	stmt, err = DB.Prepare("UPDATE userinfo SET username=? WHERE uid=?")
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getUsers() {
	rows, err := DB.Query("SELECT * FROM userinfo")
	if err != nil {
		log.Fatal(err.Error())
	}

	var uid int
	var username string
	var department string
	var created time.Time

	for rows.Next() {
		err = rows.Scan(&uid, &username, &department, &created)
		if err != nil {
			log.Fatal(err.Error())
		}

		out := fmt.Sprintf("%d %s %s %v", uid, username, department, created)
		fmt.Println(out)
	}

	rows.Close()
}

func deleteUser(uid int) {
	stmt, err := DB.Prepare("DELETE FROM userinfo WHERE uid=?")
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Deleting user:", uid)
	res, err := stmt.Exec(uid)
	if err != nil {
		log.Fatal(err.Error())
	}

	affect, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err.Error())
	}
	if affect > 0 {
		fmt.Println("Deleted user:", uid)
	} else {
		fmt.Println("No rows affected")
	}
}
