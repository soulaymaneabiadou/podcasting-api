package main

import "podcast/app"

func main() {
	app := &app.App{}

	app.Serve()
}
