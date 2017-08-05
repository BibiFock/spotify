package main

import (
	"fmt"
	"html/template"
	"net/http"
	"spotapi"
)

func main() {
	fmt.Println("I'm up get ready for this bitch")
	client := spotapi.LoadClient()

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("public/auth.html")
		t.Execute(w, client)
	})

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		client.GetToken(code)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !client.IsLogged() {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		w.Write([]byte("hello, welcome in board captain <br />" + client.GetNewSongs()))
	})

	http.ListenAndServe(":8686", nil)
}
