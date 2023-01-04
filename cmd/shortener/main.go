package main

import (
	"github.com/bbt-t/shortenerURL/configs"
	"github.com/bbt-t/shortenerURL/internal/app"
	_ "github.com/bbt-t/shortenerURL/pkg/logging"
)

func main() {
	/*
		"Application entry point".
		Parse system ENVs and start http-server.
	*/
	cfg := configs.NewConfServ()
	app.Start(cfg)
}
