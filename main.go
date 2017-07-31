package main

import (
	// "fmt"
	"html/template"
	"net/http"
	"spotapi"
)

func main() {
	client := spotapi.LoadClient()

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("public/auth.html")
		t.Execute(w, client)
	})

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		w.Write([]byte("hello, your code is now:" + code))
	})

	http.ListenAndServe(":8686", nil)
}
