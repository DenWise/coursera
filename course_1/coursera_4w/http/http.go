package main

import (
	"fmt"
	"net/http"
)

type Handler struct {
	Name string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Handler: " + h.Name, "URL: " + r.URL.String())
}

func main() {
	//rootHandler := &Handler{Name:"root"}
	//http.Handle("/", rootHandler)
	//
	//testHandler := &Handler{Name:"test"}
	//http.Handle("/test/", testHandler)
	//
	//fmt.Println("Starting server at :8082")
	//http.ListenAndServe(":8082", nil)

	//mux := http.NewServeMux()
	//mux.Handle("/", rootHandler)
	//mux.Handle("/test/", testHandler)
	//
	//server := http.Server{
	//	Addr:              ":8082",
	//	Handler:           mux,
	//	ReadTimeout:       10 * time.Second,
	//	WriteTimeout:      10 * time.Second,
	//}
	//
	//fmt.Println("Starting server at :8082")
	//server.ListenAndServe()

	go startServer(":8082")
	startServer(":8084")
}

func startServer(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Addr: " + addr, "URL: " + r.URL.String())
	})
	mux.HandleFunc("/test/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Addr: " + addr, "URL: " + r.URL.String())
	})

	server := http.Server{
		Addr:              addr,
		Handler:           mux,
	}

	fmt.Println("Starting server at ", addr)
	server.ListenAndServe()
}