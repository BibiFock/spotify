package spotapi

type Artist struct {
	Id   string
	Name string
}

func (a *Artist) GetResentSongs() {
	url := "https://api.spotify.com/v1/artists/{id}/albums"

	var sJson struct {
		Items []struct {
			Name string
			Href string
		}
		Next string
	}

	sJson.Next = strings.Replace(url, "{id}", a.Id, -1)

	fmt.Println("--- Recent songs:")
	sJson.Next = "https://api.spotify.com/v1/me/top/tracks?limit=50&time_range=long_term"
	for sJson.Next != "" {
		req, err := http.NewRequest(http.MethodGet, sJson.Next, nil)
		if err != nil {
			panic("url error")
		}

		body := c.doRequest(req)

		sJson.Next = ""
		if err := json.Unmarshal([]byte(body), &sJson); err != nil {
			panic(err.Error())
		}

		for _, item := range sJson.Items {
			fmt.Println(item)
		}

	}

}
