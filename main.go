package main

import (
	flag "github.com/spf13/pflag"
)

type Config struct {
	EnableWebServer bool
	DisableProm     bool // whether to enable prometheus metrics
	DisableWebUI    bool // whether to serve static files from Vite
	WorldDir        string
	ServerDir       string
}

var config Config

func main() {

	// web server
	flag.BoolVarP(&config.EnableWebServer, "serve", "s", false, "Start the web server")
	flag.BoolVar(&config.DisableProm, "no-prom", false, "Disable serving prometheus metrics at /api/metrics (use with --serve)")
	flag.BoolVar(&config.DisableWebUI, "no-ui", false, "Disable serving web UI (use with --serve)")

	// reading
	flag.StringVar(&config.WorldDir, "world", "./world", "Path to world directory")
	flag.StringVar(&config.ServerDir, "server", ".", "Path to server directory where usercache.json lives")

	flag.Parse()

	if config.EnableWebServer {
		initializeServer()
	} else {
		printInfo()
	}
}
