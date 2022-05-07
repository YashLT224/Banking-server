package main

import (
	"Banking/app"
	"Banking/logger"
)

func main() {

	logger.Info("starting the Application")
	app.Start()
}
