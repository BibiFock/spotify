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
	"strings"
)

const fileJson = "./client.json"
const fileTokenJson = "./client.token.json"

type Auth struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Prefix       string `json:"token_type"`
}

type Client struct {
	Url          string   `json:"url"`
	Id           string   `json:"clientId"`
	Secret       string   `json:"clientSecret"`
	ResponseType string   `json:"responseType"`
	RedirectUri  string   `json:"redirect_uri"`
	Scopes       []string `json:"scopes"`
	auth         Auth
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
	Url, err := url.Parse(api.Url)
	if err != nil {
		panic("fuck your self")
	}

	params := url.Values{}
	params.Add("client_id", api.Id)
	params.Add("response_type", api.ResponseType)
	params.Add("redirect_uri", api.RedirectUri)
	params.Add("scope", strings.Join(api.Scopes, " "))
	//params.Add("show_dialog", "true")
	Url.RawQuery = params.Encode()

	return Url.String()
}

func (c *Client) refreshToken() bool {
	urlToken := "https://accounts.spotify.com/api/token"

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

	fmt.Println(c.auth.Token)
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
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	if err := json.Unmarshal([]byte(body), &c.auth); err != nil {
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

func (c *Client) GetNewSongs() string {
	c.refreshToken()

	urlArtists := "https://api.spotify.com/v1/me/following?type=artist"
	req, err := http.NewRequest(http.MethodGet, urlArtists, nil)
	if err != nil {
		panic("url error")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", c.auth.Prefix+" "+c.auth.Token)

	spotClient := http.Client{}
	res, err := spotClient.Do(req)
	if err != nil || c.isBadResult(res) {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	return c.auth.Token + string(body)
}

type errorJson struct {
	Status  int
	Message string
}

func (c Client) isBadResult(res *http.Response) bool {
	if res.StatusCode == 200 {
		return false
	}

	if res.StatusCode == 401 {
		//fmt.Println("false false")
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err.Error())
		}

		var response map[string]errorJson
		if err := json.Unmarshal([]byte(body), &response); err != nil {
			panic(err.Error())
		}

		fmt.Println(string(body))
		if response["error"].Message == "The access token expired" {
			return !c.refreshToken()
		}

	}

	return true
}
