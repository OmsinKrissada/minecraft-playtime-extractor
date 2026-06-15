package main

import (
	"cmp"
	"embed"
	"log"
	"math"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

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

func initializeServer() {
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
	if !config.DisableWebUI {
		app.Use("/", func(c fiber.Ctx) error {
			c.Request().Header.Del("If-Modified-Since")
			return c.Next()
		})

		app.Use("/", static.New("website/dist", static.Config{
			FS: staticFiles,
		}))
	}

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

	if !config.DisableProm {
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
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
