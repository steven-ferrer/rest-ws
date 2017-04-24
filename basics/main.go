package main //every go program is made of package

//required packages are imported
//via import keyword, by convention,
//package name is the last import path
//in case of fmt, the package name is
//obviously fmt
import (
	"fmt"
	"net/http"
)

//1. Hello World
//2. Hello HTTP (simple http server)

func main() {
	http.HandleFunc("/", indexFunc)

	//Make the server listen
	http.ListenAndServe("localhost:1234", nil)
}

func indexFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello Go!")
}
