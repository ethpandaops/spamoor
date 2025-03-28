package docs

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"github.com/ethpandaops/spamoor/webui/server"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

func GetSwaggerHandler(logger logrus.FieldLogger) http.HandlerFunc {
	return httpSwagger.Handler(func(c *httpSwagger.Config) {
		c.Layout = httpSwagger.StandaloneLayout

		// override swagger header bar
		headerHTML, err := server.BuildPageHeader()
		if err != nil {
			logger.Errorf("failed generating page header for api: %v", err)
		} else {
			headerStr, err := json.Marshal(headerHTML)
			if err != nil {
				logger.Errorf("failed marshalling page header for api: %v", err)
			} else {
				var headerScript strings.Builder

				headerScript.WriteString("var headerHtml = ")
				headerScript.Write(headerStr)
				headerScript.WriteString(";")
				headerScript.WriteString("var headerEl = document.createElement(\"div\"); headerEl.className = \"header\"; headerEl.innerHTML = headerHtml; document.body.insertBefore(headerEl, document.body.firstElementChild);")
				headerScript.WriteString(`function addCss(fileName) { var el = document.createElement("link"); el.type = "text/css"; el.rel = "stylesheet"; el.href = fileName; document.head.appendChild(el); }`)
				headerScript.WriteString(`function addStyle(cssCode) { var el = document.createElement("style"); el.type = "text/css"; el.appendChild(document.createTextNode(cssCode)); document.head.appendChild(el); }`)
				headerScript.WriteString(`function addScript(fileName) { var el = document.createElement("script"); el.type = "text/javascript"; el.src = fileName; document.head.appendChild(el); }`)
				headerScript.WriteString(`addCss("/css/bootstrap.min.css");`)
				headerScript.WriteString(`addCss("/css/layout.css");`)
				headerScript.WriteString(`addScript("/js/color-modes.js");`)
				headerScript.WriteString(`addScript("/js/jquery.min.js");`)
				headerScript.WriteString(`addScript("/js/bootstrap.bundle.min.js");`)
				headerScript.WriteString(`addStyle("#swagger-ui .topbar { display: none; } .swagger-ui .opblock .opblock-section-header { background: rgba(var(--bs-body-bg-rgb), 0.8); } [data-bs-theme='dark'] .swagger-ui svg { filter: invert(100%); }");`)
				headerScript.WriteString(`
					// override swagger style (replace all color selectors)
					swaggerStyle = Array.prototype.filter.call(document.styleSheets, function(style) { return style.href && style.href.match(/swagger-ui/) })[0];
					swaggerRules = swaggerStyle.rules || swaggerStyle.cssRules;
					swaggerColorSelectors = [];
					Array.prototype.forEach.call(swaggerRules, function(rule) {
						if(rule.cssText.match(/color: rgb\(59, 65, 81\);/)) {
							swaggerColorSelectors.push(rule.selectorText);
						}
					});
					addStyle(swaggerColorSelectors.join(", ") + " { color: inherit; }");

				`)

				//nolint:gosec // ignore
				c.AfterScript = template.JS(headerScript.String())
			}
		}
	})
}
