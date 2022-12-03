package main

import (
	"flag"

	"github.com/bbt-t/shortenerURL/internal/app"
	_ "github.com/bbt-t/shortenerURL/pkg/logging"
)

func main() {
	/*
		"Application entry point".
	*/
	var inpFlagParam string
	flag.StringVar(&inpFlagParam, "db", "none", `
	"select database: no flag - use map, 'sqlite' - use SQLite, 'pg' - Postgresql, 'redis' - Redis"
	`)
	flag.Parse()

	app.Start(inpFlagParam)
}
