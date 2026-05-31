package api

import (
	"net/http"
	"strings"

	"go.mattglei.ch/go.mattglei.ch/internal/conf"
	"go.mattglei.ch/go.mattglei.ch/internal/html"
	"go.mattglei.ch/go.mattglei.ch/internal/pkg"
)

func rootEndpoint(config conf.Config, packages *pkg.Packages) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-powered-by", "go.mattglei.ch [https://github.com/gleich/go.mattglei.ch]")
		if strings.HasSuffix(r.URL.Path, "/info/refs") {
			http.Error(w, "this server does not serve git repositories", http.StatusNotFound)
			return
		}
		w.Header().Set("Cache-Control", "no-store")

		name := strings.TrimPrefix(r.URL.Path, "/")
		if name == "" {
			html.RenderIndex(config, packages, w)
			return
		}
		html.RenderPackage(config, packages, w, r)
	}
}
