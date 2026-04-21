// Package plugin provides dynamic plugin loading using Yaegi.
package plugin

import (
	"fmt"

	"github.com/ethpandaops/spamoor/utils"
	"gopkg.in/yaml.v3"
)

const (
	// PluginMetadataFile is the name of the metadata file in plugin archives.
	PluginMetadataFile = "plugin.yaml"
)

// PluginMetadata contains metadata about a plugin from plugin.yaml.
type PluginMetadata struct {
	Name       string `yaml:"name"`
	BuildTime  string `yaml:"build_time"`
	GitVersion string `yaml:"git_version"`
}

// ParsePluginMetadata parses plugin metadata from YAML bytes.
func ParsePluginMetadata(data []byte) (*PluginMetadata, error) {
	var meta PluginMetadata

	if err := yaml.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse plugin.yaml: %w", err)
	}

	if meta.Name == "" {
		return nil, fmt.Errorf("plugin.yaml: name field is required")
	}

	return &meta, nil
}

// NewLocalPluginMetadata creates metadata for a local path plugin
// using the directory name and spamoor's build info.
func NewLocalPluginMetadata(dirName string) *PluginMetadata {
	return &PluginMetadata{
		Name:       dirName,
		BuildTime:  utils.BuildTime,
		GitVersion: utils.BuildVersion,
	}
}

// GeneratePluginYAML generates the content for a plugin.yaml file.
func GeneratePluginYAML(name, buildTime, gitVersion string) ([]byte, error) {
	meta := &PluginMetadata{
		Name:       name,
		BuildTime:  buildTime,
		GitVersion: gitVersion,
	}

	return yaml.Marshal(meta)
}
