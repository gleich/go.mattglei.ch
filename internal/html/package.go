package html

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"go.mattglei.ch/go.mattglei.ch/internal/conf"
	"go.mattglei.ch/go.mattglei.ch/internal/pkg"
	"go.mattglei.ch/go.mattglei.ch/internal/util"
)

type PackageTemplate struct {
	Name string
	Root string
}

var (
	//go:embed package.html
	packageHTML     string
	packageTemplate = template.Must(template.New("gometa").Parse(packageHTML))
)

func RenderPackage(
	config conf.Config,
	packages *pkg.Packages,
	w http.ResponseWriter,
	r *http.Request,
) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	root := strings.Split(path, "/")[0]
	repo := packages.Get(root)
	if repo.Name == "" {
		http.Error(w, util.NOT_FOUND_ERROR, http.StatusNotFound)
		return
	}

	exists := repo.EnsurePath(path)
	if !exists {
		http.Error(w, util.NOT_FOUND_ERROR, http.StatusNotFound)
		return
	}

	data := PackageTemplate{
		Name: path,
		Root: root,
	}
	err := packageTemplate.Execute(w, data)
	if err != nil {
		util.InternalServerError(w, fmt.Errorf("%w failed to execute HTML template", err))
	}
}
