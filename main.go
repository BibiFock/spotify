package main

import (
	"fmt"
	"net/http"
	"spotapi"
)

func main() {
	fmt.Println("I'm up get ready for this bitch")
	client := spotapi.LoadClient()
	if !client.IsLogged() {
		srv := &http.Server{Addr: ":8686"}
		fmt.Println("Please click on this url: " + client.GetUrlAuth())

		http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
			code := r.FormValue("code")
			client.GetToken(code)
			w.Write([]byte("<html><body><script>setTimeout('window.close()', 500);</script></body>"))
			// srv.Close()
			doYourBest(client)
		})

		srv.ListenAndServe()
		// srv.Shutdown()
	} else {
		doYourBest(client)
	}
}

func doYourBest(c *spotapi.Client) {
	// c.GetNewSongs()
	fmt.Println("------")
	c.GetResentSongs()
	// c.GetLastSongs()
	// c.GetTopSongs()
}
