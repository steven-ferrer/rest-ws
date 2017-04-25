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

//DB database type, well use to interact with database
var DB *sql.DB

//User our sample resource
type User struct {
	ID         int       `json:"userid"`
	Username   string    `json:"username"`
	Department string    `json:"department"`
	Created    time.Time `json:"created"`
}

//Users collection of users
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

	log.Println("Server started on localhost:1234")

	//use the routes as multiplexer/handler
	http.ListenAndServe("localhost:1234", routes)
}

func usersGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logRequest(r)
	users, err := getUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) //500
		return
	}

	output, err := json.Marshal(users)
	if err != nil {
		//error occured while processing request
		http.Error(w, err.Error(), http.StatusInternalServerError) //500
		return
	}

	w.WriteHeader(http.StatusOK) // 200
	fmt.Fprintf(w, string(output))
}

func usersPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logRequest(r)
	username := r.FormValue("username")
	department := r.FormValue("department")

	if username == "" || department == "" {
		http.Error(w, "Cannot have empty values", http.StatusBadRequest) //400
		return
	}

	//begin creating user
	id, err := createUser(username, department, time.Now().UTC())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) //500
		return
	}

	w.WriteHeader(http.StatusOK) //200
	w.Write([]byte(fmt.Sprintf("User created. Last inserted ID is %d", id)))

}

func userGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logRequest(r)
	id := ps.ByName("id")
	if id == "" {
		log.Println("Empty ID")
		http.Error(w, "Cannot have Empty ID", http.StatusBadRequest) //400
		return
	}

	uid, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Problem reading ID", http.StatusBadRequest) //400
		return
	}

	user, err := getUser(uid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Problem getting user with ID %d", uid), http.StatusInternalServerError)
		return
	}

	output, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Problem getting user", http.StatusInternalServerError) //500
		return
	}

	w.WriteHeader(http.StatusOK)  //200
	fmt.Fprint(w, string(output)) //write to client
}

func userPut(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logRequest(r)
	//this will only modify the username, if you want,
	//you can add more fields to be updated
	id := ps.ByName("id")
	if id == "" {
		log.Println("Empty ID")
		http.Error(w, "Cannot have Empty ID's", http.StatusBadRequest) //400
		return
	}

	uid, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Problem reading ID", http.StatusBadRequest) //400
		return
	}

	username := r.FormValue("username")
	if username == "" {
		http.Error(w, "Specify a new username", http.StatusBadRequest) //400
		return
	}

	err = updateUser(uid, username)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError) //400
		return
	}

	w.WriteHeader(http.StatusOK) //200
	fmt.Fprint(w, "Update complete!")
}

func userDel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logRequest(r)
	id := ps.ByName("id")
	if id == "" {
		log.Println("Empty ID")
		http.Error(w, "Cannot have Empty ID", http.StatusBadRequest) //400
		return
	}

	uid, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Problem reading ID", http.StatusBadRequest) //400
		return
	}

	err = deleteUser(uid)
	if err != nil {
		http.Error(w, "Problem deleting user", http.StatusInternalServerError) //500
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "User deleted!")
}

//helper function to create user
func createUser(username, department string, created time.Time) (int64, error) {
	//create the new user
	stmt, err := DB.Prepare("INSERT INTO userinfo(username, departname, created) VALUES(?,?,?)")

	//fmt.Println(stmt)

	//res, err := stmt.Exec("steven", "development", "2017-04-29")
	res, err := stmt.Exec(username, department, created)
	if err != nil {
		return -1, err // -1 to indicate user was not created
	}

	//get the id of the last inserted row
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err //error getting last inserted id
	}

	fmt.Println("Last Inserted ID: ", id)
	return id, nil // no error was encountered, return the id

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
func getUsers() (Users, error) {
	//get the rows
	rows, err := DB.Query("SELECT * FROM userinfo")
	//check if an error is returned
	if err != nil {
		return Users{}, err
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
			return Users{}, err
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
	return users, nil
}

//returns a specific user
func getUser(uid int) (User, error) {
	user := User{}
	err := DB.QueryRow("SELECT * FROM userinfo WHERE uid=?", uid).Scan(&user.ID, &user.Username, &user.Department, &user.Created)
	if err != nil {
		return User{}, err
	}

	return user, nil
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

//simple logging middleware
func logRequest(r *http.Request) {
	log.Println(r.Host, r.Method, r.URL)
}
