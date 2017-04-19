package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Resp struct {
	Result int `json:"result"`
}

func addFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusForbidden)
		return
	}
	//define variables
	var n1 string = r.FormValue("num1")
	var n2 string = r.FormValue("num2")
	log.Println(n1, n2)

	if n1 == "" || n2 == "" {
		http.Error(w, "Cannot num1 and num2 should have values", http.StatusBadRequest)
		return
	}

	//shortcut for the above syntax
	num1, err := strconv.Atoi(n1)
	if err != nil {
		fmt.Fprintf(w, "Error parsing num1: %s", n1)
		return
	}
	num2, err := strconv.Atoi(n2)
	if err != nil {
		fmt.Fprintf(w, "Error parsing num2: %s", n1)
		return
	}
	res := num1 + num2

	//write to our response
	reslt := Resp{Result: res}

	out, err := json.Marshal(reslt)
	if err != nil {
		return
	}

	fmt.Fprintln(w, string(out))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Index"))
	})
	http.HandleFunc("/add", addFunc)
	//Exercise 1:
	//Implement the mult, div, and sub
	log.Println("Server listening on :9398")
	log.Fatal(http.ListenAndServe("localhost:9398", nil))
}
