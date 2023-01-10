package main

import (
	"github.com/bbt-t/shortenerURL/internal/app"
	"github.com/bbt-t/shortenerURL/internal/config"
	_ "github.com/bbt-t/shortenerURL/pkg/logging"
)

func main() {
	/*
		"Application entry point".
		Parse system ENVs and start http-server.
	*/
	cfg := config.NewConfServ()
	app.Run(cfg)
}
