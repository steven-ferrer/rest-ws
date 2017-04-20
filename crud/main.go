package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type User struct {
	ID         int       `json:"userid"`
	Username   string    `json:"username"`
	Department string    `json:"department"`
	Created    time.Time `json:"created"`
}

type Users struct {
	Users []User `json:"users"`
}

func main() {
	db, err := sql.Open("sqlite3", "../users.db")
	DB = db
	if err != nil {
		log.Fatal("Problem opening database file: ", err.Error())
	}

	//close the database when this function exits
	defer DB.Close()

	//createUser("srf", "dev", time.Now().UTC())

	//users := getUsers()
	//fmt.Println(users)

	//	for i := 1; i < 6; i++ {
	//		deleteUser(i)
	//	}

	//define the routes
	routes := httprouter.New()

	routes.GET("/users", usersGet)
	routes.POST("/users", usersPost)

	routes.GET("/users/:id", userGet)
	routes.PUT("/users/:id", userPut)
	routes.DELETE("/users/:id", userDel)
	http.ListenAndServe("localhost:1234", routes)
}

func usersGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	users := getUsers()

	output, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprintf(w, string(output))
}

func usersPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.FormValue("username")
	department := r.FormValue("department")

	if username == "" || department == "" {
		http.Error(w, "Cannot have empty values", http.StatusBadRequest)
	}

	//begin creating user
	createUser(username, department, time.Now().UTC())

}

func userGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if id == "" {
		log.Println("Empty ID")
		return
	}

	uid, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Problem reading ID", http.StatusBadRequest)
	}

	user := getUser(uid)

	output, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Problem getting user", http.StatusInternalServerError)
	}

	fmt.Fprint(w, string(output))
}

func userPut(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if id == "" {
		log.Println("Empty ID")
		return
	}

	uid, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Problem reading ID", http.StatusBadRequest)
		return
	}

	username := r.FormValue("newusername")
	if username == "" {
		http.Error(w, "Specify a new username", http.StatusBadRequest)
		return
	}

	err = updateUser(uid, username)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Update complete!")
}

func userDel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if id == "" {
		log.Println("Empty ID")
		return
	}

	uid, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Problem reading ID", http.StatusBadRequest)
		return
	}

	err = deleteUser(uid)
	if err != nil {
		http.Error(w, "Problem deleting user", http.StatusInternalServerError)
	}

	fmt.Fprint(w, "User deleted!")
}

//helper function to create user
func createUser(username, department string, created time.Time) {
	//create the new user
	stmt, err := DB.Prepare("INSERT INTO userinfo(username, departname, created) VALUES(?,?,?)")

	//fmt.Println(stmt)

	//res, err := stmt.Exec("steven", "development", "2017-04-29")
	res, err := stmt.Exec(username, department, created)
	if err != nil {
		log.Fatal(err.Error())
	}

	//get the id of the last inserted row
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Last Inserted ID: ", id)

}

//updates a user's username
func updateUser(uid int, username string) error {
	//update selected user
	_, err := DB.Exec("UPDATE userinfo SET username=? WHERE uid=?", username, uid)
	if err != nil {
		return err
	}
	return nil
}

//returns all the user from the database
func getUsers() Users {
	//get the rows
	rows, err := DB.Query("SELECT * FROM userinfo")
	//check if an error is returned
	if err != nil {
		log.Fatal(err.Error())
	}

	//close the rows when this function ends/returns
	defer rows.Close()

	//create a Users struct basically
	//an array in other languages
	users := Users{}
	for rows.Next() {
		user := User{}
		//get the contents of the current row
		err = rows.Scan(&user.ID, &user.Username, &user.Department, &user.Created)
		if err != nil {
			log.Fatal(err.Error())
		}

		//add to our array
		users.Users = append(users.Users, user)
	}

	//encode the struct to bytes
	//note the json equivalent from the above
	//output, err := json.Marshal(users)
	//if err != nil {
	//	log.Fatal(err.Error())
	//}

	//return string(output)
	return users
}

//returns a specific user
func getUser(uid int) User {
	user := User{}
	err := DB.QueryRow("SELECT * FROM userinfo WHERE uid=?", uid).Scan(&user.ID, &user.Username, &user.Department, &user.Created)
	if err != nil {
		log.Fatal(err.Error())
	}

	return user
}

func deleteUser(uid int) error {
	stmt, err := DB.Prepare("DELETE FROM userinfo WHERE uid=?")
	if err != nil {
		return err
	}

	fmt.Println("Deleting user:", uid)
	res, err := stmt.Exec(uid)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affect > 0 {
		fmt.Println("Deleted user:", uid)
	} else {
		fmt.Println("No rows affected")
	}

	return nil
}
