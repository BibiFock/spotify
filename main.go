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
		client.Token = code
		client.SaveToJson()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if client.Token == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		w.Write([]byte("hello, welcome in board captain"))
	})

	http.ListenAndServe(":8686", nil)
}
