package main

import (
	app2 "backend-processor/app"
)

func main() {
	app := app2.App{}
	app.Addr = ":1678"

	app.CertFile = "server.crt"
	app.KeyFile = "server.key"

	app.Run()
}

func Add(a int, b int) int {
	return a + b
}
