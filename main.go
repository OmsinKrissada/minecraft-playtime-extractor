package main

import (
	"fmt"
	"os"
	"path/filepath"

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

	if _, err := os.Stat(filepath.Join(config.WorldDir, "stats")); err != nil {
		fmt.Printf("Error: Cannot find stats folder in '%v'\nMake sure that '%v' contains stats directory, or use --world to change world directory.\n", config.WorldDir, config.WorldDir)
		os.Exit(1)
	}

	if _, err := os.Stat(filepath.Join(config.ServerDir, "usercache.json")); err != nil {
		fmt.Printf("Error: Cannot find usercache.json at '%v'\nMake sure that '%v' contains usercache.json file, or use --server to change server directory.\n", config.ServerDir, config.ServerDir)
		os.Exit(1)
	}

	if config.EnableWebServer {
		initializeServer()
	} else {
		printInfo()
	}
}
