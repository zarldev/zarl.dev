package main

import (
	"fmt"
	"path/filepath"

	"github.com/zarldev/zarldotdev/pkg/app"
)

func main() {
	fp := filepath.Join("config", "config.json")
	config := app.LoadConfig(fp)
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
