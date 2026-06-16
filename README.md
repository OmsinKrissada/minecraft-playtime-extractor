# Minecraft Playtime Extractor

This is a very basic tool for extracting
- Player playtime
- World run time (world ticks)
- When the player was last seen

from a Minecraft server's files.

## How It Works
This tool relies on reading game files directly. Since it's not a mod or plugin, it works on a vanilla server as well and doesn't require the Minecraft server to be running.

This table lists file/directory that this tool reads for each type of data:
|Data|Read from|
|---|---|
|World run time|`<world_dir>/level.nbt`|
|Player playtime|`<world_dir>/stats/<player_uuid>.json`|
|Player last seen timestamp|modification time of `<world_dir>/stats/<player_uuid>.json`|
|UUID -> Username mapping|`<server_dir>/usercache.json`

## Usage

When run without options, it will print out world run time and all players' playtime to standard output.

To list available options, you can run `./minecraft-playtime-extractor --help` which will give something similar to
```
Usage of ./minecraft-playtime-api:
      --no-prom         Disable serving prometheus metrics at /api/metrics (use with --serve)
      --no-ui           Disable serving web UI (use with --serve)
  -s, --serve           Start the web server
      --server string   Path to world directory (default ".")
      --world string    Path to world directory (default "./world")
```

### Web server
The default port is 8080. However you can change this by setting the "PORT" environment variable.

For example, to change the port to 3000 on a unix-like system you would run
```
PORT=3000 ./minecraft-playtime-extractor
```

There are only 2 main endpoints:
- `/api`
- `/api/metrics`: same thing as /api but for scraping by Prometheus

You can access the web UI at `/`. (http://localhost:8080/)

## Building

### Prerequisite
What you must have installed on your machine
- Node.js 22+
- Go 1.26+
- `pnpm`
- `make` (actually, you can also copy out the commands from `Makefile` to run them manually and skip installing `make` if you prefer)

### Build
Static assets for the web UI are embedded into the output binary directly. For this reason, you have to build the Vue project before compiling the Go project.

If you use wish to use the provided Makefile, you can just run
```
make
```

This will give you a binary in your project root folder.
The single binary serves as both a simple CLI tool and a web server.

## Disclaimer
The output you get is not real-time. It relies on the server to constantly write the world data to disk. This usually happens every 5 minutes on a standard vanilla server.