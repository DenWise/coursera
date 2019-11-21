package main

import (
	"fmt"
	"net/http"
	"text/template"
)

//var usersTpl = template.Must(template.ParseFiles("users.html"))

type tplParams struct {
	URL     string
	Browser string
}

type User struct {
	ID int
	Name string
	Active bool
}

func (u *User) PrintActive() string {
	if !u.Active {
		return ""
	}
	return "method says user " + u.Name + " active"
}

func IsUserOdd(u *User) bool {
	return u.ID%2 != 0
}

var users = []User{
	{
		ID:     1,
		Name:   "David",
		Active: false,
	},
	{
		ID:     2,
		Name:   "John",
		Active: true,
	},
	{
		ID:     3,
		Name:   "Derek",
		Active: true,
	},
	{
		ID:     4,
		Name:   "Robert",
		Active: false,
	},
}

const EXAMPLE = `
You use browser: {{.Browser}}

you at the page with url: {{.URL}}
`

func main() {
	http.HandleFunc("/", handle)
	http.HandleFunc("/users", fileTemplate)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
    tpl := template.New("example")
    tpl, _ = tpl.Parse(EXAMPLE)
    params := tplParams{
		URL:     r.URL.String(),
		Browser: r.UserAgent(),
	}

	tpl.Execute(w,params)
}

func fileTemplate(w http.ResponseWriter, r *http.Request) {

	tplUsersFuncMap := template.FuncMap{
		"OddUser": IsUserOdd,
	}

	tpl, err := template.
		New("").
		Funcs(tplUsersFuncMap).
		ParseFiles("users.html")
	if err != nil {
		panic(err)
	}

	err = tpl.ExecuteTemplate(w, "users.html",
		struct {
			Users []User
		}{
			users,
		})
	if err != nil {
		panic(err)
	}
	//usersTpl.Execute(w, struct {
	//	Users []User
	//}{users})
}