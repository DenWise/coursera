package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Params struct {
	ID int
	User string
}

var uploadFormTpl = []byte(`
	<html>
		<body>
			<form action="/upload" method="post" enctype="multipart/form-data">
				Image: <input type="file" name="my_file">
				<input type="submit" value="Upload">
			</form>
		</body>
	</html>
`)

func mainPage(w http.ResponseWriter, r *http.Request){
    w.Write(uploadFormTpl)
}

func uploadImage(w http.ResponseWriter, r *http.Request){
    r.ParseMultipartForm(5 * 1024 * 1025)
    file, handler, err := r.FormFile("my_file")
    if err != nil {
    	fmt.Println(err)
    	return
	}
	defer file.Close()

    fmt.Fprintf(w, "Name: %v\n", handler.Filename)
	fmt.Fprintf(w,"Header: %#v\n", handler.Header)

    hasher := md5.New()
    io.Copy(hasher, file)

	fmt.Fprintf(w,"md5: %x\n",hasher.Sum(nil))
}

func uploadRawBody(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
		http.Error(w, err.Error(),500)
    }
    defer r.Body.Close()

    p := &Params{}
    err = json.Unmarshal(body, p)
	if err != nil {
		http.Error(w, err.Error(),500)
	}

	fmt.Fprintf(w,"Content-Type: %#v\n", r.Header.Get("Content-Type"))
	fmt.Fprintf(w,"Data: %#v\n", p)
}

func main() {
	http.HandleFunc("/",mainPage)
	http.HandleFunc("/upload", uploadImage)
	http.HandleFunc("/raw_body", uploadRawBody)

	http.ListenAndServe(":8082", nil)
}