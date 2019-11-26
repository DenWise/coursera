package main

import "net/http"

func (s *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path { 
		
		case "/user/create":
			s.wrapperCreate(w, r)
		
		case "/user/profile":
			s.wrapperProfile(w, r)
		
		default:
			http.Error(w, "Method not found", 404)
		}
	}
	func (s *OtherApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path { 
		
		case "/user/create":
			s.wrapperCreate(w, r)
		
		default:
			http.Error(w, "Method not found", 404)
		}
	}