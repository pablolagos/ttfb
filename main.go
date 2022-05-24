package main

import (
	"crypto/tls"
	"embed"
	pongo22 "github.com/go-macaron/pongo2"
	"gopkg.in/macaron.v1"
	"log"
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

const CERT_KEY = "/etc/speedtest/cert.key"
const CERT_CRT = "/etc/speedtest/cert.crt"

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
	m.Get("/ttfb", controllers.CheckTTFB)

	if RunMode == "development" {
		m.Run(8080)
		os.Exit(1)
	}

	// Run server in http port
	go m.Run(80)

	tlsCert, err := tls.LoadX509KeyPair(CERT_CRT, CERT_KEY)
	if err != nil {
		log.Fatalf("Error reading ssl certificates: %s", err)
	}

	server := &http.Server{
		Addr:    "0.0.0.0:443",
		Handler: m,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
		},
	}

	log.Println("Listening TLS on " + server.Addr)
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Println(err.Error())
	}

}
