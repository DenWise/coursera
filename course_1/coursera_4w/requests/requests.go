package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func main() {
	go runServer()

	time.Sleep(100 * time.Millisecond)

	runGetData()
	runPostData()
}

func mainHandler (w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Incoming request: r %#v",r)
	fmt.Fprintf(w, "Incoming request: r.Url %#v",r.URL)
}

func postHandler (w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    defer r.Body.Close()
	if err != nil {
	    http.Error(w, err.Error(), 500)
	    return
	}
	fmt.Fprintf(w,"postHandler raw body: %s\n", string(body))
}

func getClient() *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &http.Client{
		Transport:     transport,
		Timeout:       10 * time.Second,
	}
}

func runServer() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/raw_body", postHandler)

	fmt.Println("Starting server at 8082")
	http.ListenAndServe(":8082", nil)
}

func runGetData() {
	client := getClient()

	req := &http.Request{
		Method:           http.MethodGet,
		Header:           http.Header{
			"User-Agent": {"some/user-agent"},
		},
	}

	req.URL, _ = url.Parse("http://127.0.0.1:8082/?id=43")
	req.URL.Query().Set("user", "Varien")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error happend", err)
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error happend", err)
	    return
	}

	fmt.Printf("getData: %#v\n\n\n", string(respBody))
}

func runPostData() {
	client := getClient()

	data := `{"id":42,"user":"Varien"}`
	body := bytes.NewBufferString(data)

	url := "http://127.0.0.1:8082/raw_body"
	req, _ := http.NewRequest(http.MethodPost, url, body)
	req.Header.Add("Content-Type","application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(data)))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error happend", err)
	    return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("postData %#v\n\n\n", string(respBody))
}
