package main

import "github.com/mebr0/squirrel-bot/internal/app"

const configPath = "configs/main.yml"

func main() {
	app.Run(configPath)
}
