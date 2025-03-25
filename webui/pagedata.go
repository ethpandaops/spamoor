package webui

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"syscall"
	"time"

	"github.com/ethpandaops/spamoor/utils"
	"github.com/sirupsen/logrus"
)

var LayoutTemplateFiles = []string{
	"_layout/layout.html",
	"_layout/header.html",
	"_layout/footer.html",
}

type PageData struct {
	Active         string
	Meta           *Meta
	Data           interface{}
	Version        string
	BuildTime      string
	Year           int
	Title          string
	Lang           string
	Debug          bool
	DebugTemplates []string
}

type Meta struct {
	Title       string
	Description string
	Domain      string
	Path        string
	Tlabel1     string
	Tdata1      string
	Tlabel2     string
	Tdata2      string
	Templates   string
}

type ErrorPageData struct {
	CallTime   time.Time
	CallUrl    string
	ErrorMsg   string
	StackTrace string
	Version    string
}

func InitPageData(w http.ResponseWriter, r *http.Request, active, path, title string, mainTemplates []string) *PageData {
	fullTitle := fmt.Sprintf("%v - %v - %v", title, frontendConfig.SiteName, time.Now().Year())

	if title == "" {
		fullTitle = fmt.Sprintf("%v - %v", frontendConfig.SiteName, time.Now().Year())
	}

	buildTime, _ := time.Parse("2006-01-02T15:04:05Z", utils.Buildtime)
	data := &PageData{
		Meta: &Meta{
			Title:       fullTitle,
			Description: "dugtrio: beacon chain api load balancer",
			Domain:      r.Host,
			Path:        path,
			Templates:   strings.Join(mainTemplates, ","),
		},
		Active:    active,
		Data:      &struct{}{},
		Version:   utils.GetBuildVersion(),
		BuildTime: fmt.Sprintf("%v", buildTime.Unix()),
		Year:      time.Now().UTC().Year(),
		Title:     frontendConfig.SiteName,
		Lang:      "en-US",
		Debug:     frontendConfig.Debug,
	}

	acceptedLangs := strings.Split(r.Header.Get("Accept-Language"), ",")
	if len(acceptedLangs) > 0 {
		if strings.Contains(acceptedLangs[0], "ru") || strings.Contains(acceptedLangs[0], "RU") {
			data.Lang = "ru-RU"
		}
	}

	for _, v := range r.Cookies() {
		if v.Name == "language" {
			data.Lang = v.Value
			break
		}
	}

	return data
}

// used to handle errors constructed by Template.ExecuteTemplate correctly
func HandleTemplateError(w http.ResponseWriter, r *http.Request, fileIdentifier string, functionIdentifier string, infoIdentifier string, err error) error {
	// ignore network related errors
	if err != nil && !errors.Is(err, syscall.EPIPE) && !errors.Is(err, syscall.ETIMEDOUT) {
		logger.WithFields(logrus.Fields{
			"file":       fileIdentifier,
			"function":   functionIdentifier,
			"info":       infoIdentifier,
			"error type": fmt.Sprintf("%T", err),
			"route":      r.URL.String(),
		}).WithError(err).Error("error executing template")
		http.Error(w, "Internal server error", http.StatusServiceUnavailable)
	}
	return err
}

func HandlePageError(w http.ResponseWriter, r *http.Request, pageError error) {
	templateFiles := append(LayoutTemplateFiles, "_layout/500.html")
	notFoundTemplate := GetTemplate(templateFiles...)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusInternalServerError)
	data := InitPageData(w, r, "blockchain", r.URL.Path, "Internal Error", templateFiles)
	errData := &ErrorPageData{
		CallTime: time.Now(),
		CallUrl:  r.URL.String(),
		ErrorMsg: pageError.Error(),
		Version:  utils.GetBuildVersion(),
	}
	data.Data = errData
	err := notFoundTemplate.ExecuteTemplate(w, "layout", data)
	if err != nil {
		logrus.Errorf("error executing page error template for %v route: %v", r.URL.String(), err)
		http.Error(w, "Internal server error", http.StatusServiceUnavailable)
	}
}

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	templateFiles := append(LayoutTemplateFiles, "_layout/404.html")
	notFoundTemplate := GetTemplate(templateFiles...)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	data := InitPageData(w, r, "blockchain", r.URL.Path, "Not Found", templateFiles)
	err := notFoundTemplate.ExecuteTemplate(w, "layout", data)
	if err != nil {
		logrus.Errorf("error executing not-found template for %v route: %v", r.URL.String(), err)
		http.Error(w, "Internal server error", http.StatusServiceUnavailable)
	}
}
