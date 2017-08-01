package spotapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

const fileJson = "./client.json"

type Client struct {
	Url          string   `json:"url"`
	Id           string   `json:"clientId"`
	Secret       string   `json:"clientSecret"`
	ResponseType string   `json:"responseType"`
	RedirectUri  string   `json:"redirect_uri"`
	Scopes       []string `json:"scopes"`
	Token        string   `json:"token"`
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
	Url.RawQuery = params.Encode()

	return Url.String()
}

func (c Client) SaveToJson() {
	bytes, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if ioutil.WriteFile(fileJson, bytes, 0644) != nil {
		panic("can't write")
	}

}
