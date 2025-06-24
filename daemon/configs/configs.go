package configs

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
)

// SpammerConfig represents a spammer configuration for export/import.
// This uses the same format as StartupSpammerConfig to maintain compatibility.
type SpammerConfig struct {
	Scenario    string                 `yaml:"scenario"`
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	Config      map[string]interface{} `yaml:"config"`
}

// ConfigImportItem represents either a spammer config or an include directive
type ConfigImportItem struct {
	// Spammer configuration fields
	Scenario    string                 `yaml:"scenario,omitempty"`
	Name        string                 `yaml:"name,omitempty"`
	Description string                 `yaml:"description,omitempty"`
	Config      map[string]interface{} `yaml:"config,omitempty"`

	// Include directive
	Include string `yaml:"include,omitempty"`
}

// ResolveConfigImports recursively resolves includes and returns the final spammer configs
func ResolveConfigImports(input string, baseURL string, visited map[string]bool) ([]SpammerConfig, error) {
	// Resolve the actual source path/URL
	resolvedInput := resolveConfigIncludePath(input, baseURL)

	// Prevent circular includes
	if visited[resolvedInput] {
		return nil, fmt.Errorf("circular include detected: %s", resolvedInput)
	}
	visited[resolvedInput] = true
	defer func() { delete(visited, resolvedInput) }()

	// Get the YAML data
	var yamlData string
	var err error

	// Check if resolved input is a URL
	if isConfigURL(resolvedInput) {
		yamlData, err = downloadConfigFromURL(resolvedInput)
		if err != nil {
			return nil, fmt.Errorf("failed to download from URL: %w", err)
		}
	} else if isConfigFilePath(resolvedInput) {
		// Try to read as file path
		yamlData, err = readConfigFromFile(resolvedInput)
		if err != nil {
			return nil, fmt.Errorf("failed to read from file: %w", err)
		}
	} else {
		// Treat as raw YAML data (only for root input)
		yamlData = resolvedInput
	}

	// Parse as ConfigImportItems to handle both configs and includes
	var importItems []ConfigImportItem
	if err := yaml.Unmarshal([]byte(yamlData), &importItems); err != nil {
		return nil, fmt.Errorf("failed to parse import data: %w", err)
	}

	var allConfigs []SpammerConfig

	// Determine the new base URL/path for nested includes
	newBaseURL := getConfigBaseURL(resolvedInput)

	// Process each item
	for _, item := range importItems {
		if item.Include != "" {
			// This is an include directive
			includedConfigs, err := ResolveConfigImports(item.Include, newBaseURL, visited)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve include '%s': %w", item.Include, err)
			}
			allConfigs = append(allConfigs, includedConfigs...)
		} else {
			// This is a spammer configuration
			config := SpammerConfig{
				Scenario:    item.Scenario,
				Name:        item.Name,
				Description: item.Description,
				Config:      item.Config,
			}
			allConfigs = append(allConfigs, config)
		}
	}

	return allConfigs, nil
}

// MergeScenarioConfiguration merges scenario defaults with provided configuration
func MergeScenarioConfiguration(scenario *scenario.Descriptor, providedConfig map[string]interface{}) (string, error) {
	// Get default configurations
	defaultYaml, err := yaml.Marshal(scenario.DefaultOptions)
	if err != nil {
		return "", fmt.Errorf("failed to marshal default config: %w", err)
	}

	defaultWalletConfig := spamoor.GetDefaultWalletConfig(scenario.Name)
	defaultWalletConfigYaml, err := yaml.Marshal(defaultWalletConfig)
	if err != nil {
		return "", fmt.Errorf("failed to marshal default wallet config: %w", err)
	}

	// Merge configurations
	mergedConfig := map[string]interface{}{}

	if err := yaml.Unmarshal(defaultWalletConfigYaml, &mergedConfig); err != nil {
		return "", fmt.Errorf("failed to unmarshal default wallet config: %w", err)
	}

	if err := yaml.Unmarshal(defaultYaml, &mergedConfig); err != nil {
		return "", fmt.Errorf("failed to unmarshal default config: %w", err)
	}

	// Apply provided config overrides
	for k, v := range providedConfig {
		mergedConfig[k] = v
	}

	configYAML, err := yaml.Marshal(mergedConfig)
	if err != nil {
		return "", fmt.Errorf("failed to marshal merged config: %w", err)
	}

	return string(configYAML), nil
}

// isConfigURL checks if the input string is a valid URL
func isConfigURL(input string) bool {
	parsedURL, err := url.Parse(input)
	if err != nil {
		return false
	}
	return parsedURL.Scheme == "http" || parsedURL.Scheme == "https"
}

// isConfigFilePath checks if the input string looks like a file path
func isConfigFilePath(input string) bool {
	// Check if it contains YAML content markers (likely raw YAML)
	if len(input) > 0 && (input[0] == '-' || input[0] == '[' || input[0] == '{') {
		return false
	}

	// Check if file exists
	if _, err := os.Stat(input); err == nil {
		return true
	}

	// Check if it looks like a path (contains / or \)
	return len(input) > 0 && (input[0] == '/' || input[0] == '.' || input[0] == '~' ||
		(len(input) > 1 && input[1] == ':')) // Windows drive letter
}

// readConfigFromFile reads YAML data from a local file
func readConfigFromFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	return string(data), nil
}

// downloadConfigFromURL downloads YAML data from a remote URL
func downloadConfigFromURL(urlStr string) (string, error) {
	// Validate URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
	}

	// Download the YAML data
	resp, err := http.Get(urlStr)
	if err != nil {
		return "", fmt.Errorf("failed to download from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error %d: %s", resp.StatusCode, resp.Status)
	}

	yamlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(yamlData), nil
}

// resolveConfigIncludePath resolves an include path against a base URL or directory
func resolveConfigIncludePath(includePath, baseURL string) string {
	// If include path is absolute (URL or absolute file path), return as-is
	if isConfigURL(includePath) || filepath.IsAbs(includePath) {
		return includePath
	}

	// If no base URL provided, return the include path as-is
	if baseURL == "" {
		return includePath
	}

	// If base is a URL, resolve relative URL
	if isConfigURL(baseURL) {
		baseURLParsed, err := url.Parse(baseURL)
		if err != nil {
			return includePath // fallback to original
		}

		resolvedURL, err := baseURLParsed.Parse(includePath)
		if err != nil {
			return includePath // fallback to original
		}

		return resolvedURL.String()
	}

	// If base is a file path, resolve relative to directory
	if isConfigFilePath(baseURL) {
		baseDir := filepath.Dir(baseURL)
		return filepath.Join(baseDir, includePath)
	}

	// Fallback to original include path
	return includePath
}

// getConfigBaseURL extracts the base URL or directory from a source path
func getConfigBaseURL(sourcePath string) string {
	if isConfigURL(sourcePath) {
		// For URLs, get the base URL (everything except the filename)
		parsedURL, err := url.Parse(sourcePath)
		if err != nil {
			return ""
		}

		// Remove the filename from the path
		parsedURL.Path = path.Dir(parsedURL.Path)
		if !strings.HasSuffix(parsedURL.Path, "/") {
			parsedURL.Path += "/"
		}

		return parsedURL.String()
	}

	if isConfigFilePath(sourcePath) {
		// For file paths, return the directory
		return filepath.Dir(sourcePath)
	}

	// For raw YAML data, no base URL
	return ""
}
