package main

import (
	flag "github.com/spf13/pflag"
)

type Config struct {
	EnableWebServer bool
	DisableProm     bool // whether to enable prometheus metrics
	DisableWebUI    bool // whether to serve static files from Vite
}

func main() {
	var config Config

	// web server
	flag.BoolVarP(&config.EnableWebServer, "serve", "s", false, "Start the web server")
	flag.BoolVar(&config.DisableProm, "no-prom", false, "Disable serving prometheus metrics at /api/metrics (use with --serve)")
	flag.BoolVar(&config.DisableWebUI, "no-ui", false, "Disable serving web UI (use with --serve)")

	flag.Parse()

	if config.EnableWebServer {
		initializeServer(!config.DisableWebUI, !config.DisableProm)
	} else {
		printInfo()
	}
}
