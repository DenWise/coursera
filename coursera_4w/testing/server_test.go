package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type TestCaseCheckout struct {
	ID string
	Result *CheckoutResult
	IsError bool
}

func TestCartCheckout(t *testing.T) {
	cases := []TestCaseCheckout{
		{
			ID:      "42",
			Result:  &CheckoutResult{
				Status:  200,
				Balance: 100500,
				Err:     "",
			},
			IsError: false,
		},
		{
			ID:      "100500",
			Result:  &CheckoutResult{
				Status:  400,
				Balance: 0,
				Err:     "bad_balance",
			},
			IsError: false,
		},
		{
			ID:      "__broken_json",
			Result:  nil,
			IsError: true,
		},
		{
			ID:      "__internal_error",
			Result:  nil,
			IsError: true,
		},
	}

	// create new test server
	ts := httptest.NewServer(http.HandlerFunc(CheckoutDummy))

	for caseNum, item := range cases {
		c := &Cart{PaymentApiURL:ts.URL}
		result, err := c.Checkout(item.ID)

		if err != nil && !item.IsError {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && item.IsError {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !reflect.DeepEqual(result, item.Result) {
			t.Errorf("[%d] wrong result: got %#v, expected %#v", caseNum, result, item.Result)
		}
	}
	ts.Close()
}
