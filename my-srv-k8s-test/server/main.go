package main 

import (
	"fmt"
	"net/http"
)

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page!")
}

func hwPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hw")
}

func main(){
	http.HandleFunc("/", homePageHandler)
	http.HandleFunc("/hw", hwPageHandler)

	http.ListenAndServe(":8080", nil)
	fmt.Print("server is running on http://localhost:8080")
}