package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

func main() {
	api := getApiInfos()
	fmt.Println(api.getUrlAuth())
}

func getApiInfos() ApiClient {
	raw, err := ioutil.ReadFile("./client.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var api ApiClient

	if err := json.Unmarshal([]byte(raw), &api); err != nil {
		panic("boom")
	}

	return api
}

type ApiClient struct {
	Url          string   `json:"url"`
	Id           string   `json:"clientId"`
	Secret       string   `json:"clientSecret"`
	ResponseType string   `json:"responseType"`
	RedirectUri  string   `json:"redirect_uri"`
	Scopes       []string `json:"scopes"`
}

func (api ApiClient) toString() string {
	return toJson(api)
}

func (api ApiClient) getUrlAuth() string {
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

func toJson(api interface{}) string {
	bytes, err := json.Marshal(api)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(bytes)
}
