package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w,"Hello, I'm your first http response")
	w.Write([]byte("Horay!!!"))
}

func main() {
	http.HandleFunc("/page",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Single page: ", r.URL.String())
		})

	http.HandleFunc("/pages/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Multiple pages: ", r.URL.String())
		})

	http.HandleFunc("/",handler)

	fmt.Println("Starting server at :8082")
	http.ListenAndServe(":8082", nil)
}