package main

import (
	"fmt"
	"spotapi"
)

func main() {
	api := spotapi.LoadClient()
	fmt.Println(api.GetUrlAuth())
}
