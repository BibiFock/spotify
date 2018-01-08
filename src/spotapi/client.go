package spotapi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const fileJson = "./client.json"
const fileTokenJson = "./client.token.json"

type Auth struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Prefix       string `json:"token_type"`
}

type Client struct {
	Id          string   `json:"clientId"`
	Secret      string   `json:"clientSecret"`
	RedirectUri string   `json:"redirect_uri"`
	scopes      []string `json:"scopes"`
	auth        Auth
	artists     []Artist
}

func LoadClient() (c *Client) {
	raw, err := ioutil.ReadFile(fileJson)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := json.Unmarshal([]byte(raw), &c); err != nil {
		panic("boom")
	}

	raw, err = ioutil.ReadFile(fileTokenJson)
	if err != nil {
		fmt.Println(fileTokenJson + " not found")
		return c
	}

	if err = json.Unmarshal([]byte(raw), &c.auth); err != nil {
		panic("boom 2 ")
	}

	return c
}

func (api Client) GetUrlAuth() string {
	var Url *url.URL
	Url, err := url.Parse("https://accounts.spotify.com/authorize/")
	if err != nil {
		panic("fuck your self")
	}

	scopes := []string{
		"playlist-read-private",
		"playlist-read-collaborative",
		"playlist-modify-public",
		"playlist-modify-private",
		"streaming",
		"ugc-image-upload",
		"user-follow-modify",
		"user-follow-read",
		"user-library-read",
		"user-library-modify",
		"user-read-private",
		"user-read-birthdate",
		"user-read-email",
		"user-top-read",
		"user-read-recently-played",
	}

	params := url.Values{}
	params.Add("client_id", api.Id)
	params.Add("response_type", "code")
	params.Add("redirect_uri", api.RedirectUri)
	params.Add("scope", strings.Join(scopes, " "))
	Url.RawQuery = params.Encode()

	return Url.String()
}

func (c *Client) refreshToken() bool {
	urlToken := "https://accounts.spotify.com/api/token"
	// fmt.Println("old token: " + c.auth.Token)

	params := url.Values{}
	params.Add("grant_type", "refresh_token")
	params.Add("refresh_token", c.auth.RefreshToken)

	req, err := http.NewRequest(http.MethodPost, urlToken, strings.NewReader(params.Encode()))
	if err != nil {
		panic("url error")
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(c.Id + ":" + c.Secret))
	req.Header.Add("Authorization", "Basic "+encoded)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	spotClient := http.Client{}
	res, err := spotClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	if res.StatusCode != 200 {
		fmt.Println(" refresh request response: " + string(body))
		return false
	}

	if err := json.Unmarshal([]byte(body), &c.auth); err != nil {
		panic(err)
	}

	// fmt.Println("new token: " + c.auth.Token)
	// fmt.Println("-----------")
	c.saveCurrentToken()
	return true
}

func (c *Client) GetToken(code string) {
	urlToken := "https://accounts.spotify.com/api/token"

	params := url.Values{}
	params.Add("grant_type", "authorization_code")
	params.Add("code", code)
	params.Add("redirect_uri", c.RedirectUri)
	params.Add("client_id", c.Id)
	params.Add("client_secret", c.Secret)

	req, err := http.NewRequest(http.MethodPost, urlToken, strings.NewReader(params.Encode()))
	if err != nil {
		panic("url error")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	spotClient := http.Client{}
	res, err := spotClient.Do(req)
	if err != nil {
		fmt.Println("fail fatal")
		log.Fatal(err)
	}

	// fmt.Println(res.Header)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("read fail fail")
		panic(err.Error())
	}

	if err := json.Unmarshal([]byte(body), &c.auth); err != nil {
		fmt.Println("unmarshal fail fail")
		panic(err)
	}

	c.saveCurrentToken()
}

func (c Client) saveCurrentToken() {
	bytes, err := json.Marshal(c.auth)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if ioutil.WriteFile(fileTokenJson, bytes, 0644) != nil {
		panic("can't write")
	}
}

func (c Client) IsLogged() bool {
	return c.auth != (Auth{})
}

func (c *Client) GetLastSongs() {
	var sJson struct {
		Items []struct {
			Track struct {
				Artists []Artist
				Name    string
			}
			Played_at string
		}
		Next string
	}

	fmt.Println("--- Last songs:")
	sJson.Next = "https://api.spotify.com/v1/me/player/recently-played?limit=50"
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

		fmt.Println("---> test list ")
		for _, item := range sJson.Items {
			fmt.Println(item)
		}
		fmt.Println(sJson.Next)
	}
}

func (c *Client) GetTopSongs() {
	var sJson struct {
		Items []struct {
			Name string
		}
		Next string
	}

	fmt.Println("--- Top songs:")
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

func (c *Client) GetFollowingNewSongs() {
	c.loadFollowingArtists()
	var sAlbums struct {
		Next  string
		Items []struct {
			// Album_type    string
			External_urls struct {
				Spotify string
			}
			Id   string
			Name string
			// Type string
		}
	}

	var albumUrl = "https://api.spotify.com/v1/albums/"
	var sAlbum struct {
		Name                   string
		Release_date           string
		Release_date_precision string
	}

	for _, artist := range c.artists {
		sAlbums.Next = "https://api.spotify.com/v1/artists/" + artist.Id + "/albums?market=FR&limit=50&album_type=album"
		for sAlbums.Next != "" {
			req, err := http.NewRequest(http.MethodGet, sAlbums.Next, nil)
			if err != nil {
				panic("url error")
			}

			body := c.doRequest(req)

			sAlbums.Next = ""
			if err := json.Unmarshal([]byte(body), &sAlbums); err != nil {
				panic(err.Error())
			}

			for _, album := range sAlbums.Items {
				req, err := http.NewRequest(http.MethodGet, albumUrl+album.Id, nil)
				if err != nil {
					panic("url error")
				}

				body := c.doRequest(req)

				if err := json.Unmarshal([]byte(body), &sAlbum); err != nil {
					panic(err.Error())
				}

				if sAlbum.Release_date_precision != "day" {
					continue
				}

				date, err := time.Parse("2006-01-02", sAlbum.Release_date)
				if err != nil {
					fmt.Println("can't compare "+sAlbum.Release_date, err)
				}
				if time.Since(date).Hours() < 7*4*24 {
					fmt.Println(artist.Name + " - " + album.Name + " (" + sAlbum.Release_date + ")")
					break
				}
			}
		}
	}
}

func (c *Client) GetNewSongs() {
	c.loadFollowingArtists()
	var sJson struct {
		Albums struct {
			Next  string
			Items []struct {
				Artists       []Artist
				Name          string
				Id            string
				External_urls struct {
					Spotify string
				}
			}
		}
	}

	fmt.Println("--- New Songs are:")
	sJson.Albums.Next = "https://api.spotify.com/v1/browse/new-releases?limit=50"
	for sJson.Albums.Next != "" {
		// fmt.Println(sJson.Albums.Next)
		req, err := http.NewRequest(http.MethodGet, sJson.Albums.Next, nil)
		if err != nil {
			panic("url error")
		}

		body := c.doRequest(req)

		sJson.Albums.Next = ""
		if err := json.Unmarshal([]byte(body), &sJson); err != nil {
			panic(err.Error())
		}

		for _, album := range sJson.Albums.Items {
			for _, art := range album.Artists {
				found := false
				for _, artist := range c.artists {
					if art.Id != artist.Id {
						continue
					}
					found = true
				}
				if found {
					fmt.Println(art.Name + " - " + album.Name + " (" + album.External_urls.Spotify + ")")
					break
				}
			}
		}
	}
}

func (c *Client) loadFollowingArtists() {
	var jsonStruct struct {
		Artists struct {
			Next  string
			Items []Artist
		}
	}

	jsonStruct.Artists.Next = "https://api.spotify.com/v1/me/following?type=artist&limit=50"

	for jsonStruct.Artists.Next != "" {
		req, err := http.NewRequest(http.MethodGet, jsonStruct.Artists.Next, nil)
		if err != nil {
			panic("url error")
		}

		body := c.doRequest(req)

		jsonStruct.Artists.Next = ""
		if err := json.Unmarshal([]byte(body), &jsonStruct); err != nil {
			panic(err.Error())
		}

		for _, artist := range jsonStruct.Artists.Items {
			c.artists = append(c.artists, artist)
		}
	}
}

func (c *Client) doRequest(req *http.Request) string {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", c.auth.Prefix+" "+c.auth.Token)

	spotClient := http.Client{}
	res, err := spotClient.Do(req)
	if res.StatusCode == 429 {
		second, _ := strconv.Atoi(res.Header.Get("Retry-After"))
		time.Sleep(time.Duration(second+1) * time.Second)
		res, err = spotClient.Do(req)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	if res.StatusCode == 200 {
		return string(body)
	}

	var response map[string]struct {
		Status  int
		Message string
	}

	if err := json.Unmarshal([]byte(body), &response); err != nil {
		panic(err.Error())
	}

	if res.StatusCode == 401 {
		if response["error"].Message == "The access token expired" {
			fmt.Println("Need to refresh fucking token")
			if !c.refreshToken() {
				panic("failed to refresh token")
			}

			return c.doRequest(req)
		}
	}

	fmt.Println(strconv.Itoa(res.StatusCode) + ": " + response["error"].Message)
	return strconv.Itoa(res.StatusCode) + ": " + response["error"].Message
}
