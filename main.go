package main

import (
	"embed"
	pongo22 "github.com/go-macaron/pongo2"
	"gopkg.in/macaron.v1"
	"net/http"
	"os"
	"speedtest/binfiles"
	"speedtest/controllers"
	"time"
)

//go:embed public
var staticFS embed.FS

//go:embed templates
var teplatesFS embed.FS

var RunMode = "development"

func main() {

	m := macaron.New()
	macaron.Env = RunMode
	m.Use(macaron.Static("public",
		macaron.StaticOptions{
			FileSystem:  binfiles.New(&staticFS, "public"),
			SkipLogging: true,
			ETag:        true,
			Expires:     func() string { return time.Now().Add(time.Minute * time.Duration(60)).Format(http.TimeFormat) },
		},
	))
	m.Use(macaron.Recovery())

	m.Use(pongo22.Pongoer(pongo22.Options{
		TemplateFileSystem: binfiles.New(&teplatesFS, "templates"),
	}))

	m.Get("/", controllers.Home)
	m.Post("/ttfb", controllers.CheckTTFB)

	if RunMode == "development" {
		m.Run(8080)
	} else {
		m.Run(80)
	}
	os.Exit(0)

}
