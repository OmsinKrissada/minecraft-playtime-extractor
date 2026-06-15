package main

import (
	"cmp"
	"compress/gzip"
	"embed"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/Tnze/go-mc/nbt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const SERVER_DIR = "."
const WORLD_DIR_NAME = "world"

const TICKS_IN_AN_HOUR = 20 * 60 * 60

type UserCache struct {
	Name string
	Uuid string
}

func getUsernameMap() map[string]string {
	path := filepath.Join(SERVER_DIR, "usercache.json")
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var decoded []UserCache
	json.Unmarshal(data, &decoded)

	uuid_username_map := make(map[string]string)
	for _, uc := range decoded {
		uuid_username_map[uc.Uuid] = uc.Name
	}

	return uuid_username_map
}

func getSinglePlaytime(path string) int {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	type DecodedType struct {
		Stats struct {
			Custom struct {
				Playtime int `json:"minecraft:play_time"`
			} `json:"minecraft:custom"`
		}
	}
	var decoded DecodedType
	json.Unmarshal(data, &decoded)
	pt := decoded.Stats.Custom.Playtime

	return pt
}

func getSingleLastSeen(path string) time.Time {
	filestat, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return filestat.ModTime()
}

func getAllPlaytime() map[string]int {
	path := filepath.Join(SERVER_DIR, WORLD_DIR_NAME, "stats")
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	user_playtime_tick := make(map[string]int)
	for _, f := range files {
		playerStatPath := filepath.Join(path, f.Name())
		pt := getSinglePlaytime(playerStatPath)
		username, _, _ := strings.Cut(f.Name(), ".")
		user_playtime_tick[username] = pt
	}

	return user_playtime_tick
}

func getAllLastSeen() map[string]time.Time {
	path := filepath.Join(SERVER_DIR, WORLD_DIR_NAME, "stats")
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	player_lastseen := make(map[string]time.Time)
	for _, f := range files {
		playerStatPath := filepath.Join(path, f.Name())
		ls := getSingleLastSeen(playerStatPath)
		username, _, _ := strings.Cut(f.Name(), ".")
		player_lastseen[username] = ls
	}

	return player_lastseen
}

type WorldNBT struct {
	Data struct {
		Time int64
	}
}

func getWorldRunTime() int64 {
	path := filepath.Join(SERVER_DIR, WORLD_DIR_NAME, "level.dat")
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}
	defer gzipReader.Close()

	var decoded WorldNBT
	nbtDecoder := nbt.NewDecoder(gzipReader)
	nbtDecoder.Decode(&decoded)
	return decoded.Data.Time
}

type APIResponse struct {
	Name       string  `json:"name"`
	Uuid       string  `json:"uuid"`
	PlaytimeHr float64 `json:"playtime_hr"`
	LastSeen   string  `json:"last_seen"`
}

func transformResponse(pt_map map[string]int, ls_map map[string]time.Time) []APIResponse {
	var arr []APIResponse

	uc := getUsernameMap()
	for uuid, pt := range pt_map {
		pt_hours := float64(pt) / TICKS_IN_AN_HOUR
		pt_hours = math.Round(pt_hours*100) / 100
		arr = append(arr, APIResponse{Name: uc[uuid], Uuid: uuid, PlaytimeHr: pt_hours, LastSeen: ls_map[uuid].Format(time.RFC3339)})
	}

	slices.SortFunc(arr, func(a, b APIResponse) int {
		return cmp.Compare(b.PlaytimeHr, a.PlaytimeHr)
	})

	return arr
}

//go:embed all:website/dist
var staticFiles embed.FS

func main() {
	app := fiber.New(fiber.Config{
		TrustProxy:  true,
		ProxyHeader: fiber.HeaderXForwardedFor,
		TrustProxyConfig: fiber.TrustProxyConfig{
			Private: true,
		},
	})

	app.Use(logger.New())
	app.Use(cors.New())

	// embed FS doesn't store modtime properly
	// these 4 lines are for preventing permanent cache due to 304 status code on static files
	app.Use("/", func(c fiber.Ctx) error {
		c.Request().Header.Del("If-Modified-Since")
		return c.Next()
	})

	app.Use("/", static.New("website/dist", static.Config{
		FS: staticFiles,
	}))

	app.Get("/api", func(c fiber.Ctx) error {
		pt := getAllPlaytime()
		ls := getAllLastSeen()
		transformed := transformResponse(pt, ls)
		worldRunTime := float64(getWorldRunTime()) / TICKS_IN_AN_HOUR
		worldRunTime = math.Round(worldRunTime*100) / 100

		resp := map[string]any{
			"players":        transformed,
			"world_run_time": worldRunTime,
		}
		return c.JSON(resp)
	})

	app.Get("/api/metrics", adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reg := prometheus.NewRegistry()
		playtime := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "minecraft_star_tech_playtime_hours",
			},
			[]string{"name", "uuid"},
		)
		reg.MustRegister(playtime)

		pt := getAllPlaytime()
		ls := getAllLastSeen()
		transformed := transformResponse(pt, ls)

		for _, p := range transformed {
			labels := prometheus.Labels{"name": p.Name, "uuid": p.Uuid}
			playtime.With(labels).Set(p.PlaytimeHr)
		}

		promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(w, r)
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
