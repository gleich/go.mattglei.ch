package html

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"

	"go.mattglei.ch/go.mattglei.ch/internal/conf"
	"go.mattglei.ch/go.mattglei.ch/internal/github"
	"go.mattglei.ch/go.mattglei.ch/internal/pkg"
	"go.mattglei.ch/go.mattglei.ch/internal/util"
)

type IndexTemplate struct {
	Config   conf.Config
	Packages []github.Repository
}

var (
	//go:embed index.html
	indexHTML     string
	indexTemplate = template.Must(template.New("gometa").Parse(indexHTML))
)

func RenderIndex(config conf.Config, packages *pkg.Packages, w http.ResponseWriter) {
	data := IndexTemplate{
		Packages: packages.All(),
		Config:   config,
	}
	err := indexTemplate.Execute(w, data)
	if err != nil {
		util.InternalServerError(w, fmt.Errorf("executing html template: %w", err))
	}
}
