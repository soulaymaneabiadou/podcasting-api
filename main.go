package main

import (
	"os"
	"podcast/app"
)

func main() {
	args := os.Args[1:]

	app := &app.App{}

	if len(args) > 0 {
		arg := os.Args[1]
		if arg == "db:seed" {
			app.Seed()
		}
		if arg == "db:migrate" {
			app.Migrate()
		}
	} else {
		app.Serve()
	}
}
