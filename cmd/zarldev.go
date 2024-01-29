package main

import (
	"fmt"

	"github.com/zarldev/zarldotdev/pkg/app"
)

func main() {
	config := app.LoadConfig("./config.json")
	app, err := app.New(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = app.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
}
