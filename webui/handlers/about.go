package handlers

import (
	"net/http"
	"runtime/debug"
	"sort"

	"github.com/ethpandaops/spamoor/utils"
	"github.com/ethpandaops/spamoor/webui/server"
)

type AboutPage struct {
	Repository string
	License    string
	GoVersion  string
	Version    string
	BuildInfo  *debug.BuildInfo
	Modules    []ModuleInfo
}

type ModuleInfo struct {
	Path        string
	Version     string
	Replacement string
}

// About will return the "about" page using a go template
func (fh *FrontendHandler) About(w http.ResponseWriter, r *http.Request) {
	var templateFiles = append(server.LayoutTemplateFiles,
		"about/about.html",
	)

	var pageTemplate = server.GetTemplate(templateFiles...)
	data := server.InitPageData(r, "about", "/about", "About", templateFiles)

	var pageError error
	data.Data, pageError = fh.getAboutPageData()
	if pageError != nil {
		server.HandlePageError(w, r, pageError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if server.HandleTemplateError(w, r, "about.go", "About", "", pageTemplate.ExecuteTemplate(w, "layout", data)) != nil {
		return
	}
}

func (fh *FrontendHandler) getAboutPageData() (*AboutPage, error) {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		buildInfo = &debug.BuildInfo{}
	}

	modules := make([]ModuleInfo, 0, len(buildInfo.Deps))
	for _, dep := range buildInfo.Deps {
		module := ModuleInfo{
			Path:    dep.Path,
			Version: dep.Version,
		}

		if dep.Replace != nil {
			module.Replacement = dep.Replace.Path + " " + dep.Replace.Version
		}

		modules = append(modules, module)
	}

	// Sort modules by path
	sort.Slice(modules, func(i, j int) bool {
		return modules[i].Path < modules[j].Path
	})

	return &AboutPage{
		Repository: "https://github.com/ethpandaops/spamoor",
		License:    "MIT",
		GoVersion:  buildInfo.GoVersion,
		Version:    utils.GetBuildVersion(),
		BuildInfo:  buildInfo,
		Modules:    modules,
	}, nil
}
