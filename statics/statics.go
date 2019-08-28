// Package statics allows items to be served from domain.com/item
package statics

import (
	"fmt"
	"net/http"

	"github.com/aeridya/core"
	"github.com/aeridya/core/configurit"
	"github.com/aeridya/core/logit"
	"github.com/aeridya/module"
)

// Statics is an object that embeds Aeridya's Handler into it to allow for a list of
// handler functions to process each static file.
// Dir is the directory where statics are on the file system
// Statics.Statics is a slice of strings to serve via this system
type statics struct {
	module.Module
	Dir     string
	Statics []string
}

var (
	stats *statics
)

func init() {
	stats = &statics{Statics: make([]string, 0)}
}

// Init will create the Statics Object and return it
// Requires a path to the statics dir
func Init(opts ...module.Option) error {
	stats.ParseOpts(opts)
	return nil
}

// ReadConfig is a function to read the configuratoin file
// for all settings
func FromConfig() (err error) {
	if stats.Dir, err = configurit.Config.GetString("statics", "directory"); err != nil {
		return fmt.Errorf("Unable to get Statics from config: %s", err)
	}
	return nil
}

// Defaults places the recommended default items in the Statics handler
// Meaning:  favicon.ico, sitemap.xml, and robots.txt
// Warning:  Items must be placed in the root of the Statics Directory
func Defaults() module.Option {
	return func() {
		stats.Statics = append(stats.Statics, "/favicon.ico")
		stats.Statics = append(stats.Statics, "/sitemap.xml")
		stats.Statics = append(stats.Statics, "/robots.txt")
	}
}

func Favicons() module.Option {
	return func() {
		stats.Statics = append(stats.Statics, "/apple-icon-57x57.png")
		stats.Statics = append(stats.Statics, "/apple-icon-60x60.png")
		stats.Statics = append(stats.Statics, "/apple-icon-72x72.png")
		stats.Statics = append(stats.Statics, "/apple-icon-76x76.png")
		stats.Statics = append(stats.Statics, "/apple-icon-114x114.png")
		stats.Statics = append(stats.Statics, "/apple-icon-120x120.png")
		stats.Statics = append(stats.Statics, "/apple-icon-144x144.png")
		stats.Statics = append(stats.Statics, "/apple-icon-152x152.png")
		stats.Statics = append(stats.Statics, "/apple-icon-180x180.png")
		stats.Statics = append(stats.Statics, "/android-icon-192x192.png")
		stats.Statics = append(stats.Statics, "/favicon-32x32.png")
		stats.Statics = append(stats.Statics, "/favicon-96x96.png")
		stats.Statics = append(stats.Statics, "/manifest.json")
		stats.Statics = append(stats.Statics, "/ms-icon-144x144.png")
	}
}

func Directory(s string) module.Option {
	return func() {
		stats.Dir = s
	}
}

// Add requires an item, adds item into Statics array
func Add(item string) {
	stats.Statics = append(stats.Statics, item)
}

// serve creates a http.Handler and adds it to the DefaultServeMux of the application
func serve(pattern string, filename string) {
	http.Handle(pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	}))
}

// Serve iterates over Statics array and creates a handler for each static object
// Requires array of http.Handlers for Aeridya Handler compatibility
func Start() {
	for t := range stats.Statics {
		if core.Development {
			logit.Logf(logit.DEBUG, "Adding Static Handler for %s to directory %s", stats.Statics[t], stats.Dir+stats.Statics[t])
		}
		serve(stats.Statics[t], stats.Dir+stats.Statics[t])
	}
}
