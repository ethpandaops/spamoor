package handlers

import (
	"bytes"
	"html/template"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ethpandaops/spamoor/webui/server"
)

func TestPprofLoginTemplateRenders(t *testing.T) {
	base := "../templates"
	files := make([]string, 0, len(server.LayoutTemplateFiles)+1)
	for _, f := range server.LayoutTemplateFiles {
		files = append(files, filepath.Join(base, f))
	}
	files = append(files, filepath.Join(base, "debug/pprof_login.html"))

	tmpl, err := template.New("layout").Funcs(server.GetTemplateFuncs()).ParseFiles(files...)
	if err != nil {
		t.Fatalf("failed to parse templates: %v", err)
	}

	var out bytes.Buffer
	if err := tmpl.ExecuteTemplate(&out, "page", &struct{}{}); err != nil {
		t.Fatalf("failed to execute page: %v", err)
	}
	if !strings.Contains(out.String(), "pprof-login-required") {
		t.Errorf("page is missing the login prompt")
	}

	out.Reset()
	if err := tmpl.ExecuteTemplate(&out, "js", &struct{}{}); err != nil {
		t.Fatalf("failed to execute js: %v", err)
	}
	// The cookie name is shared with the server side gate.
	if !strings.Contains(out.String(), "spamoor_debug_token") {
		t.Errorf("js does not set the debug token cookie")
	}
}
