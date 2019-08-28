package trailslash

import (
	"net/http"

	"github.com/aeridya/core"
	"github.com/aeridya/module"
)

type trailslash struct {
	module.Module
	trail bool
}

var (
	opt *trailslash
)

func init() {
	opt = &trailslash{}
}

func Set(options ...module.Option) {
	opt.ParseOpts(options)
	if opt.trail {
		core.AddHandler(999, AddTrailingSlash)
	} else {
		core.AddHandler(999, NoTrailingSlash)
	}
}

func Slash() module.Option {
	return func() {
		opt.trail = true
	}
}

func NoSlash() module.Option {
	return func() {
		opt.trail = false
	}
}

func NoTrailingSlash(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			u := r.URL.Path
			if len(u) > 1 {
				if u[len(u)-1:] == "/" {
					u = u[:len(u)-1]
					http.Redirect(w, r, u, 301)
					return
				}
			}
		}
		h.ServeHTTP(w, r)
	})
}

func AddTrailingSlash(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			u := r.URL.Path
			if len(u) > 1 {
				if u[len(u)-1:] != "/" {
					u = u + "/"
					http.Redirect(w, r, u, 301)
					return
				}
			}
		}
		h.ServeHTTP(w, r)
	})
}
