package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type Cart struct {
	PaymentApiURL string
}

type CheckoutResult struct {
	Status int
	Balance int
	Err string
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("id")
	if key == "42" {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"status": 200, "resp": {"user": 42}}`)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status": 500, "err": "db_error"}`)
	}
}

func (c *Cart) Checkout(id string) (*CheckoutResult, error) {
	url := c.PaymentApiURL + "?id=" + id
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &CheckoutResult{}

	err = json.Unmarshal(data, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// here we mocking the remote server handling our request
func CheckoutDummy(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("id") // server get the id
	switch key {
	case "42": // in case when checkout is OK
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"status": 200, "balance": 100500}`)
	case "100500": // in case when there is not enough money on balance
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"status": 400, "err": "bad_balance"}`)
	case "__broken_json": // in case when json is malformed
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"status": 400`) //broken json
	case "__internal_error": // in case when internal server error
		fallthrough
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
