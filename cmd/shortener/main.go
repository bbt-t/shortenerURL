package main

import (
	"flag"

	"github.com/bbt-t/shortenerURL/internal/app"
	_ "github.com/bbt-t/shortenerURL/pkg/logging"
)

func main() {
	/*
		"Application entry point".
		Parse launch parameters and start http-server.
	*/
	var inpFlagParam string

	flag.StringVar(&inpFlagParam, "db", "none", `
	"select database: no flag - use map, 'sqlite' - use SQLite, 'pg' - Postgresql, 'redis' - Redis"
	`)
	flag.Parse()
	// pkg.StopNotifyAdmin()
	app.Start(inpFlagParam /* received flag or nothing ("") */)
}
