package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// SpammerConfig represents a single spammer configuration
type SpammerConfig struct {
	Scenario    string                 `yaml:"scenario"`
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	Config      map[string]interface{} `yaml:"config"`
}

// IndexEntry represents a spammer config in the index
type IndexEntry struct {
	File         string   `yaml:"file"`
	Name         string   `yaml:"name"`
	Description  string   `yaml:"description"`
	Tags         []string `yaml:"tags"`
	SpammerCount int      `yaml:"spammer_count"`
	Scenarios    []string `yaml:"scenarios"`
	MinVersion   string   `yaml:"min_version,omitempty"`
}


// Index represents the complete spammer index
type Index struct {
	Generated time.Time    `yaml:"generated"`
	Configs   []IndexEntry `yaml:"configs"`
}

// HeaderInfo contains metadata extracted from file headers
type HeaderInfo struct {
	Name        string
	Description string
	Tags        string
	MinVersion  string
}

func main() {
	configDir := "spammer-configs"
	if len(os.Args) > 1 {
		configDir = os.Args[1]
	}

	index, err := generateIndex(configDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating index: %v\n", err)
		os.Exit(1)
	}

	outputFile := filepath.Join(configDir, "_index.yaml")
	err = writeIndex(index, outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing index: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated spammer index with %d configs\n", len(index.Configs))
}

func generateIndex(configDir string) (*Index, error) {
	files, err := filepath.Glob(filepath.Join(configDir, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to list YAML files: %w", err)
	}

	// Filter out the index file itself
	var configFiles []string
	for _, file := range files {
		if !strings.HasSuffix(file, "_index.yaml") {
			configFiles = append(configFiles, file)
		}
	}

	index := &Index{
		Generated: time.Now().UTC(),
		Configs:   []IndexEntry{},
	}

	for _, file := range configFiles {
		err := processConfigFile(file, index)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to process %s: %v\n", file, err)
			continue
		}
	}

	return index, nil
}

func processConfigFile(filePath string, index *Index) error {
	// Extract header metadata
	headerInfo, err := extractHeaderInfo(filePath)
	if err != nil {
		return fmt.Errorf("failed to extract header info: %w", err)
	}

	// Parse YAML content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var configs []SpammerConfig
	err = yaml.Unmarshal(data, &configs)
	if err != nil {
		// Try parsing as single config
		var singleConfig SpammerConfig
		err = yaml.Unmarshal(data, &singleConfig)
		if err != nil {
			return fmt.Errorf("failed to parse YAML: %w", err)
		}
		configs = []SpammerConfig{singleConfig}
	}

	if len(configs) == 0 {
		return fmt.Errorf("no spammer configs found in file")
	}

	// Parse tags from comma-separated string - pass through exactly as defined
	var tags []string
	if headerInfo.Tags != "" {
		for _, tag := range strings.Split(headerInfo.Tags, ",") {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				tags = append(tags, tag)
			}
		}
	}
	
	// Collect all unique scenarios
	scenarioMap := make(map[string]bool)
	for _, config := range configs {
		if config.Scenario != "" {
			scenarioMap[config.Scenario] = true
		}
	}
	
	var scenarios []string
	for scenario := range scenarioMap {
		scenarios = append(scenarios, scenario)
	}
	sort.Strings(scenarios)

	// Get filename from path
	fileName := filepath.Base(filePath)
	
	// Use header name or first config name
	name := headerInfo.Name
	if name == "" {
		name = configs[0].Name
	}

	// Use header description or first config description
	description := headerInfo.Description
	if description == "" {
		description = configs[0].Description
	}

	// Create index entry
	entry := IndexEntry{
		File:         fileName,
		Name:         name,
		Description:  description,
		Tags:         tags,
		SpammerCount: len(configs),
		Scenarios:    scenarios,
		MinVersion:   headerInfo.MinVersion,
	}

	// Add to flat list
	index.Configs = append(index.Configs, entry)

	return nil
}

func extractHeaderInfo(filePath string) (*HeaderInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info := &HeaderInfo{}
	scanner := bufio.NewScanner(file)
	
	// Look for header comments at the beginning of the file
	commentRegex := regexp.MustCompile(`^#\s*(\w+):\s*(.+)`)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Stop processing if we hit non-comment content
		if line != "" && !strings.HasPrefix(line, "#") {
			break
		}
		
		matches := commentRegex.FindStringSubmatch(line)
		if len(matches) == 3 {
			key := strings.ToLower(matches[1])
			value := strings.TrimSpace(matches[2])
			
			switch key {
			case "name":
				info.Name = value
			case "description":
				info.Description = value
			case "tags":
				info.Tags = value
			case "min_version":
				info.MinVersion = value
			}
		}
	}

	return info, scanner.Err()
}


func writeIndex(index *Index, outputFile string) error {
	// Sort configs by name for consistent output
	sort.Slice(index.Configs, func(i, j int) bool {
		return index.Configs[i].Name < index.Configs[j].Name
	})

	data, err := yaml.Marshal(index)
	if err != nil {
		return fmt.Errorf("failed to marshal index: %w", err)
	}

	// Add header comment
	header := `# Auto-generated spammer configuration index
# Generated: ` + index.Generated.Format(time.RFC3339) + `
# DO NOT EDIT MANUALLY - This file is automatically generated by scripts/generate-spammer-index.go

`

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(header)
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write index data: %w", err)
	}

	return nil
}