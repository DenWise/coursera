package main

import (
	"encoding/json"
	"encoding/xml"
	"github.com/gorilla/schema"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

var decoder = schema.NewDecoder()

var (
	orderFields = [3]string{"Id","Age","Name"}
)

type XMLRoot struct {
	XMLName xml.Name `xml:"root"`
	Rows []XMLUser `xml:"row"`
}

type XMLUser struct {
	XMLName xml.Name `xml:"row"`
	Id int `xml:"id"`
	FirstName string `xml:"first_name"`
	LastName string `xml:"last_name"`
	Age int `xml:"age"`
	About string `xml:"about"`
	Gender string `xml:"gender"`
	Name string `xml:"-"`
}

type TestAccessCase struct {
	Client *SearchClient
	StatusCode int
	Error bool
}

type TestLimitOffsetCase struct {
	Request SearchRequest
	Error bool
}

type TestBadServerCase struct {
	Handler http.HandlerFunc
	Error bool
}

type TestFindUserCase struct {
	Request SearchRequest
	Response *SearchResponse
	Error bool
}

func TestAccessToken(t *testing.T) {
	cases := []TestAccessCase{
		{
			Client: &SearchClient{
				AccessToken: "777",
			},
			Error: false,
		},
		{
			Client: &SearchClient{
				AccessToken: "111",
			},
			Error: true,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	for caseNum, item := range cases {
		item.Client.URL = ts.URL
		searchReq := SearchRequest{
			Limit:      5,
			Offset:     0,
			Query:      "test",
			OrderField: "Name",
			OrderBy:    0,
		}
		_, err := item.Client.FindUsers(searchReq)

		if err == nil && item.Error {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
	}
	ts.Close()
}

func TestLimitOffset(t *testing.T) {
	cases := []TestLimitOffsetCase{
		{
			Request: SearchRequest{
				Limit:      10,
				Offset:     5,
			},
			Error:   false,
		},
		{
			Request: SearchRequest{
				Limit:      26,
				Offset:     5,
			},
			Error:   false,
		},
		{
			Request: SearchRequest{
				Limit:      -2,
				Offset:     5,
			},
			Error:   true,
		},
		{
			Request: SearchRequest{
				Limit:      10,
				Offset:     -7,
			},
			Error:   true,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	for caseNum, item := range cases {
		client := &SearchClient{
			AccessToken: "777",
			URL:         ts.URL,
		}
		_, err := client.FindUsers(item.Request)

		if err == nil && item.Error {
			t.Errorf("[%d] expected error got nil", caseNum)
		}
	}

	ts.Close()
}

func TestBadServer(t *testing.T) {
	cases := []TestBadServerCase{
		{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(11 * time.Second)
				w.WriteHeader(http.StatusGatewayTimeout)
			}),
			Error:   true,
		},
		{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}),
			Error:   true,
		},
		{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			}),
			Error:   true,
		},
		{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				io.WriteString(w,`{"error":"ErrorBadOrderField"}`)
			}),
			Error:   true,
		},
		{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				io.WriteString(w,`{"error":"UnknownError"}`)
			}),
			Error:   true,
		},
		{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w,"Error",5)
			}),
			Error:   true,
		},
		{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			Error:   false,
		},
	}

	for caseNum, item := range cases {
		ts := httptest.NewServer(item.Handler)
		searchReq := SearchRequest{
			Limit:      10,
			Offset:     5,
		}
		client := &SearchClient{
			AccessToken: "777",
			URL:         ts.URL,
		}
		_, err := client.FindUsers(searchReq)

		if err == nil && item.Error {
			t.Errorf("[%d] expected error got nil", caseNum)
		}

		ts.Close()
	}
}

func TestFindUsers(t *testing.T) {
	cases := []TestFindUserCase{
		{
			Request:  SearchRequest{
				Limit:      5,
				Offset:     0,
				Query:      "may",
				OrderField: "Name",
				OrderBy:    0,
			},
			Response: &SearchResponse{
				Users:    []User{
					{
						Id:     1,
						Name:   "Hilda Mayer",
						Age:    21,
						About:  "Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.\n",
						Gender: "female",
					},
					{
						Id:     6,
						Name:   "Jennings Mays",
						Age:    39,
						About:  "Veniam consectetur non non aliquip exercitation quis qui. Aliquip duis ut ad commodo consequat ipsum cupidatat id anim voluptate deserunt enim laboris. Sunt nostrud voluptate do est tempor esse anim pariatur. Ea do amet Lorem in mollit ipsum irure Lorem exercitation. Exercitation deserunt adipisicing nulla aute ex amet sint tempor incididunt magna. Quis et consectetur dolor nulla reprehenderit culpa laboris voluptate ut mollit. Qui ipsum nisi ullamco sit exercitation nisi magna fugiat anim consectetur officia.\n",
						Gender: "male",
					},
				},
				NextPage: false,
			},
			Error: false,
		},
		{
			Request:  SearchRequest{
				Limit:      2,
				Offset:     0,
				Query:      "may",
				OrderField: "Name",
				OrderBy:    0,
			},
			Response: &SearchResponse{
				Users:    []User{
					{
						Id:     1,
						Name:   "Hilda Mayer",
						Age:    21,
						About:  "Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.\n",
						Gender: "female",
					},
					{
						Id:     6,
						Name:   "Jennings Mays",
						Age:    39,
						About:  "Veniam consectetur non non aliquip exercitation quis qui. Aliquip duis ut ad commodo consequat ipsum cupidatat id anim voluptate deserunt enim laboris. Sunt nostrud voluptate do est tempor esse anim pariatur. Ea do amet Lorem in mollit ipsum irure Lorem exercitation. Exercitation deserunt adipisicing nulla aute ex amet sint tempor incididunt magna. Quis et consectetur dolor nulla reprehenderit culpa laboris voluptate ut mollit. Qui ipsum nisi ullamco sit exercitation nisi magna fugiat anim consectetur officia.\n",
						Gender: "male",
					},
				},
				NextPage: false,
			},
			Error: false,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	for caseNum, item := range cases {

		client := &SearchClient{
			AccessToken: "777",
			URL:         ts.URL,
		}

		result, err := client.FindUsers(item.Request)

		if err == nil && item.Error {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if err != nil && !item.Error {
			t.Errorf("[%d] unexpected error %s", caseNum, err)
		}
		if !reflect.DeepEqual(result, item.Response) {
			t.Errorf("[%d] wrong results: got %#v, expected %#v",
				caseNum, result, item.Response)
		}
	}
	ts.Close()
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	var (
		request SearchRequest
		orderField = "Name"
		founded = make([]XMLUser,0,100)
	)

	incomingToken := r.Header.Get("AccessToken")

	if incomingToken != "777" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := decoder.Decode(&request, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if request.OrderField != "" && !inArray(request.OrderField,orderFields) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(w, `{"error":"ErrorBadOrderField"}`)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	} else {
		orderField = request.OrderField
	}

	users := getXMLUsers()
	query := strings.ToLower(request.Query)

	for _, u := range users {
		if strings.Contains(strings.ToLower(u.Name),query) || strings.Contains(strings.ToLower(u.About),query) {
			founded = append(founded,u)
		}
	}

	if orderField == "Id" {
		sort.SliceStable(founded, func(i, j int) bool {
			return founded[i].Id < founded[j].Id
		})
	} else if orderField == "Age" {
		sort.SliceStable(founded, func(i, j int) bool {
			return founded[i].Age < founded[j].Age
		})
	} else {
		sort.SliceStable(founded, func(i, j int) bool {
			return founded[i].Name < founded[j].Name
		})
	}

	responseJson, err := json.Marshal(founded)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = io.WriteString(w,string(responseJson))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getXMLUsers() []XMLUser {
	xmlFile, err := os.Open("dataset.xml")
	if err != nil {
		log.Fatal("error while opening data xml file",err)
	}
	defer xmlFile.Close()

	xmlAsByteArray, _ := ioutil.ReadAll(xmlFile)

	var root XMLRoot
	err = xml.Unmarshal(xmlAsByteArray, &root)
	if err != nil {
		log.Fatal("cannot unmarshal data xml", err)
	}

	for i, user := range root.Rows {
		root.Rows[i].Name = user.FirstName + " " + user.LastName
	}

	return root.Rows
}

func inArray(val interface{}, array interface{}) (exists bool) {
	exists = false

	switch reflect.TypeOf(array).Kind() {
	case reflect.Array:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				exists = true
				return
			}
		}
	}

	return
}