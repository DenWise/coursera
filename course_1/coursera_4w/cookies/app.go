package main

import (
	"fmt"
	"net/http"
	"time"
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	loggedIn := (err != http.ErrNoCookie)

	if loggedIn {
		fmt.Fprintln(w,`<a href="/logout">Logout</a>`)
		fmt.Fprintln(w,"Hello, ", session.Value)
	} else {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w,`<a href="/login">Login</a></br><img src="/data/img/gopher.png" />`)
		fmt.Fprintln(w,"You need to login")
	}
}

func loginPage(w http.ResponseWriter, r *http.Request){
	expiration := time.Now().Add(10 * time.Hour)
	cookie := http.Cookie{
		Name:       "session_id",
		Value:      "Varien",
		Expires:    expiration,
	}

	http.SetCookie(w,&cookie)
	http.Redirect(w,r,"/",http.StatusFound)
}

func logoutPage(w http.ResponseWriter, r *http.Request){
    session, err := r.Cookie("session_id")
    if err == http.ErrNoCookie {
    	http.Redirect(w,r,"/",http.StatusFound)
	}

	session.Expires = time.Now().AddDate(0,0,-1)
	http.SetCookie(w, session)

	http.Redirect(w,r,"/",http.StatusFound)
}

func main() {
	// handler for static resources
	staticHandler := http.StripPrefix("/data/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/logout", logoutPage)
	http.HandleFunc("/", mainPage)
	http.Handle("/data/", staticHandler)

	fmt.Println("Starting server at :8082")
	http.ListenAndServe(":8082",nil)
}