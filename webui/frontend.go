package webui

import (
	"embed"
	"errors"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/ethpandaops/spamoor/webui/types"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New().WithField("module", "frontend")
var frontendConfig *types.FrontendConfig

var (
	//go:embed static/*
	staticEmbedFS embed.FS

	//go:embed templates/*
	templateEmbedFS embed.FS
)

type Frontend struct {
	defaultHandler  http.Handler
	rootFileSys     http.FileSystem
	router          *mux.Router
	NotFoundHandler func(http.ResponseWriter, *http.Request)
}

func NewFrontend(config *types.FrontendConfig) (*Frontend, error) {
	frontendConfig = config

	subFs, err := fs.Sub(staticEmbedFS, "static")
	if err != nil {
		return nil, err
	}
	fileSys := http.FS(subFs)

	templateFiles, err = fs.Sub(templateEmbedFS, "templates")
	if err != nil {
		return nil, err
	}

	frontend := Frontend{
		defaultHandler:  http.FileServer(fileSys),
		rootFileSys:     fileSys,
		NotFoundHandler: HandleNotFound,
	}

	return &frontend, nil
}

func (frontend *Frontend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// basically a copy of http.FileServer and of the first lines http.serveFile functions
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	name := path.Clean(upath)
	f, err := frontend.rootFileSys.Open(name)
	if err != nil {
		handleHTTPError(err, frontend.NotFoundHandler, w, r)
		return
	}
	defer f.Close()

	_, err = f.Stat()
	if err != nil {
		handleHTTPError(err, frontend.NotFoundHandler, w, r)
		return
	}

	if strings.HasSuffix(r.URL.Path, "/") {
		handleHTTPError(fs.ErrNotExist, frontend.NotFoundHandler, w, r)
		return
	}

	frontend.defaultHandler.ServeHTTP(w, r)
}

func handleHTTPError(err error, handler func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request) {
	// If error is 404, use custom handler
	if errors.Is(err, fs.ErrNotExist) {
		handler(w, r)
		return
	}
	// otherwise serve http error
	if errors.Is(err, fs.ErrPermission) {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return
	}
	// Default:
	logrus.WithError(err).Errorf("page handler error")
	http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
}
