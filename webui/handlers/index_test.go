package handlers

import (
	"html/template"
	"io"
	"path/filepath"
	"testing"

	"github.com/ethpandaops/spamoor/webui/server"
)

func TestIndexTemplateParses(t *testing.T) {
	base := "../templates"
	files := make([]string, 0, len(server.LayoutTemplateFiles)+1)
	for _, f := range server.LayoutTemplateFiles {
		files = append(files, filepath.Join(base, f))
	}
	files = append(files, filepath.Join(base, "index/index.html"))

	tmpl, err := template.New("layout").Funcs(server.GetTemplateFuncs()).ParseFiles(files...)
	if err != nil {
		t.Fatalf("failed to parse templates: %v", err)
	}

	group := &IndexPageSpammer{
		ID: 1, Name: "g", Description: "d", IsGroup: true,
		ThroughputMode: "shared", TotalThroughput: 10,
		AutoRestartFailed: true, AutoRestartCooldown: 300,
		Members: []*IndexPageSpammer{},
	}
	if err := tmpl.ExecuteTemplate(io.Discard, "groupHeaderRow", group); err != nil {
		t.Fatalf("failed to execute groupHeaderRow: %v", err)
	}
}
