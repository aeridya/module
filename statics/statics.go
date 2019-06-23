// Package statics allows items to be served from domain.com/item
package statics

import (
	"net/http"

	"github.com/aeridya/core"
	"github.com/aeridya/module"
)

// Statics is an object that embeds Aeridya's Handler into it to allow for a list of
// handler functions to process each static file.
// Dir is the directory where statics are on the file system
// Statics.Statics is a slice of strings to serve via this system
type statics struct {
	module.Module
	core.Handler
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

	return nil
}

// Defaults places the recommended default items in the Statics handler
// Meaning:  favicon.ico, sitemap.xml, and robots.txt
// Warning:  Items must be placed in the root of the Statics Directory
func Defaults() {
	stats.Statics = append(stats.Statics, "/favicon.ico")
	stats.Statics = append(stats.Statics, "/sitemap.xml")
	stats.Statics = append(stats.Statics, "/robots.txt")
}

// Add requires an item, adds item into Statics array
func Add(item string) {
	stats.Statics = append(stats.Statics, item)
}

// serve creates a http.Handler and adds it to the DefaultServeMux of the application
func serve(pattern string, filename string) {
	http.Handle(pattern, stats.Final(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})))
}

// Serve iterates over Statics array and creates a handler for each static object
// Requires array of http.Handlers for Aeridya Handler compatibility
func (s Statics) Serve(handles []func(http.Handler) http.Handler) {
	s.AppendHandlers(handles)
	for t := range s.Statics {
		s.serve(s.Statics[t], s.Dir+s.Statics[t])
	}
}
