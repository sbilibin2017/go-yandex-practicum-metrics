package main

import (
	"go-yandex-practicum-metrics/cmd/server/app"
	"os"
)

func main() {
	cmd := app.NewCommand()
	code := app.Run(cmd)
	os.Exit(code)
}
