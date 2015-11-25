package main

import (
	"flag"
	"github.com/codegangsta/martini"
	"log"
	"net/http"
	// "github.com/codegangsta/martini-contrib/gzip"
	"github.com/curt-labs/HitchInfo/controllers"
	"os"
)

var (
	listenAddr = flag.String("http", ":3000", "http listen address")
	port       = os.Getenv("PORT")
)

var m *martini.Martini

/**
 * All GET routes require either public or private api keys to be passed in.
 *
 * All POST routes require private api keys to be passed in.
 */
func main() {
	flag.Parse()

	m = martini.New()

	// os.Setenv("PORT", *listenAddr)
	if port == "" {
		port = *listenAddr
	}

	// Setup Middleware
	m.Use(martini.Recovery())
	m.Use(martini.Logger())
	// m.Use(gzip.All())
	m.Use(martini.Static("static"))

	r := martini.NewRouter()

	// r.Get("/replace", controllers.MassReplace)
	r.Get("/index.cfm", controllers.IndexRedirect)

	r.Get("/(.*)", controllers.Index)

	m.Action(r.Handle)

	log.Println(http.ListenAndServe(port, m))
}
