// Package plugin provides dynamic plugin loading using Yaegi.
package plugin

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/ethpandaops/spamoor/plugin/symbols"
	"github.com/ethpandaops/spamoor/scenario"
	"github.com/sirupsen/logrus"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

const (
	// PluginBasePath is the base path where plugins are mounted in the virtual GOPATH.
	PluginBasePath = "src/github.com/ethpandaops/spamoor/plugins"
)

// PluginLoader handles loading plugins from tar archives at runtime.
type PluginLoader struct {
	logger *logrus.Entry
}

// NewPluginLoader creates a new plugin loader.
func NewPluginLoader(logger logrus.FieldLogger) *PluginLoader {
	return &PluginLoader{
		logger: logger.WithField("component", "plugin-loader"),
	}
}

// LoadFromBytes loads a plugin from tar(.gz) bytes.
// The pluginName is used to construct the proper GOPATH directory structure.
// If compressed is true, the data is treated as gzip-compressed.
func (l *PluginLoader) LoadFromBytes(pluginName string, data []byte, compressed bool) (*Descriptor, error) {
	return l.LoadFromReader(pluginName, bytes.NewReader(data), compressed)
}

// LoadFromFile loads a plugin from a tar(.gz) file path.
// The plugin name is derived from the filename (without extension).
// Compression is auto-detected based on the .gz extension.
func (l *PluginLoader) LoadFromFile(filePath string) (*Descriptor, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin file: %w", err)
	}
	defer file.Close()

	// Determine plugin name and compression from filename
	baseName := filepath.Base(filePath)
	compressed := strings.HasSuffix(baseName, ".gz")

	// Strip extensions to get plugin name
	pluginName := baseName
	if compressed {
		pluginName = strings.TrimSuffix(pluginName, ".gz")
	}

	pluginName = strings.TrimSuffix(pluginName, ".tar")

	return l.LoadFromReader(pluginName, file, compressed)
}

// LoadFromReader loads a plugin from a tar(.gz) stream.
// The pluginName is used to construct the proper GOPATH directory structure.
// If compressed is true, the stream is treated as gzip-compressed.
func (l *PluginLoader) LoadFromReader(pluginName string, data io.Reader, compressed bool) (*Descriptor, error) {
	var tarReader io.Reader = data

	if compressed {
		gzReader, err := gzip.NewReader(data)
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzReader.Close()

		tarReader = gzReader
	}

	// Build in-memory filesystem from tar
	memFS, err := l.buildMemoryFS(tarReader, pluginName)
	if err != nil {
		return nil, fmt.Errorf("failed to build memory filesystem: %w", err)
	}

	return l.loadFromFS(pluginName, memFS)
}

// buildMemoryFS creates an in-memory filesystem from tar content.
// The files are placed under the plugin path structure.
func (l *PluginLoader) buildMemoryFS(tarReader io.Reader, pluginName string) (fs.FS, error) {
	memFS := newMemoryFS()
	tr := tar.NewReader(tarReader)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("failed to read tar entry: %w", err)
		}

		// Skip directories, we create them implicitly
		if header.Typeflag == tar.TypeDir {
			continue
		}

		// Only process regular files
		if header.Typeflag != tar.TypeReg {
			continue
		}

		// Construct the full path within the virtual GOPATH
		// Structure: src/github.com/ethpandaops/spamoor/plugins/<plugin_name>/<file_path>
		cleanPath := filepath.Clean(header.Name)
		// Remove leading "./" or "/" if present
		cleanPath = strings.TrimPrefix(cleanPath, "./")
		cleanPath = strings.TrimPrefix(cleanPath, "/")

		fullPath := filepath.Join(PluginBasePath, pluginName, cleanPath)

		// Read file content
		content, err := io.ReadAll(tr)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", header.Name, err)
		}

		memFS.addFile(fullPath, content, header.FileInfo().Mode())
	}

	return memFS, nil
}

// loadFromFS interprets a plugin from a filesystem.
func (l *PluginLoader) loadFromFS(pluginName string, filesys fs.FS) (*Descriptor, error) {
	// Create interpreter with the filesystem as GOPATH source
	i := l.newInterpreter(filesys)

	// Import the main plugin package
	pluginPkg := fmt.Sprintf("github.com/ethpandaops/spamoor/plugins/%s", pluginName)

	// First, import the package
	_, err := i.Eval(fmt.Sprintf(`import plugin "%s"`, pluginPkg))
	if err != nil {
		return nil, l.wrapEvalError(pluginName, err)
	}

	// Then evaluate the PluginDescriptor
	v, err := i.Eval("plugin.PluginDescriptor")
	if err != nil {
		return nil, l.wrapEvalError(pluginName, err)
	}

	// Extract the descriptor
	desc, ok := extractPluginDescriptor(v)
	if !ok {
		return nil, fmt.Errorf("PluginDescriptor in %s is not of type plugin.Descriptor (got %T)", pluginName, v.Interface())
	}

	l.logger.Infof("loaded plugin: %s with %d scenarios", desc.Name, len(desc.Scenarios))

	return desc, nil
}

// newInterpreter creates a new Yaegi interpreter with filesystem support.
func (l *PluginLoader) newInterpreter(filesys fs.FS) *interp.Interpreter {
	i := interp.New(interp.Options{
		GoPath:               ".",
		SourcecodeFilesystem: filesys,
	})

	// Load standard library symbols
	if err := i.Use(stdlib.Symbols); err != nil {
		l.logger.Warnf("failed to load stdlib symbols: %v", err)
	}

	// Load spamoor package symbols (generated by yaegi extract)
	if err := i.Use(symbols.Symbols); err != nil {
		l.logger.Warnf("failed to load spamoor symbols: %v", err)
	}

	return i
}

// wrapEvalError wraps Yaegi evaluation errors with helpful hints.
func (l *PluginLoader) wrapEvalError(pluginName string, err error) error {
	errStr := err.Error()

	// Try to detect missing symbols and suggest yaegi extract
	if strings.Contains(errStr, "undefined:") {
		// Parse error like: undefined: "github.com/some/package".Function
		re := regexp.MustCompile(`undefined: "([^"]+)"\.(\w+)`)
		matches := re.FindStringSubmatch(errStr)
		if len(matches) == 3 {
			pkg := matches[1]
			symbol := matches[2]
			return fmt.Errorf("failed to evaluate plugin %s: missing symbol %s.%s\n"+
				"Hint: The package %q may not have extracted symbols.\n"+
				"Run: yaegi extract %s",
				pluginName, pkg, symbol, pkg, pkg)
		}

		// Simpler undefined pattern
		simpleRe := regexp.MustCompile(`undefined: (\w+)`)
		simpleMatches := simpleRe.FindStringSubmatch(errStr)
		if len(simpleMatches) == 2 {
			return fmt.Errorf("failed to evaluate plugin %s: %w (hint: symbol %q is undefined, check imports and ensure all required packages have extracted symbols)",
				pluginName, err, simpleMatches[1])
		}
	}

	// Detect CFG/panic errors (Yaegi internal issues)
	if strings.Contains(errStr, "CFG") || strings.Contains(errStr, "panic") {
		return fmt.Errorf("failed to evaluate plugin %s: %w\n"+
			"Hint: This may be a Yaegi interpreter limitation.\n"+
			"Common causes:\n"+
			"  - Channel sends to custom type aliases in closures (use raw channel type)\n"+
			"  - Complex type assertions\n"+
			"  - Unsupported language features",
			pluginName, err)
	}

	// Detect type-related errors
	if strings.Contains(errStr, "type") && (strings.Contains(errStr, "cannot") || strings.Contains(errStr, "mismatch")) {
		return fmt.Errorf("failed to evaluate plugin %s: %w (hint: type mismatch detected, ensure types match the extracted symbols exactly)",
			pluginName, err)
	}

	// Detect import errors
	if strings.Contains(errStr, "import") || strings.Contains(errStr, "could not import") {
		return fmt.Errorf("failed to evaluate plugin %s: %w (hint: import failed, check that all imported packages have been extracted with 'yaegi extract')",
			pluginName, err)
	}

	// Generic error
	return fmt.Errorf("failed to evaluate plugin %s: %w", pluginName, err)
}

// extractPluginDescriptor attempts to extract a plugin.Descriptor from a reflect.Value.
func extractPluginDescriptor(v reflect.Value) (*Descriptor, bool) {
	if !v.IsValid() {
		return nil, false
	}

	// Handle pointer types
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	// Try direct type assertion on the interface
	iface := v.Interface()

	// Check if it's already a *Descriptor
	if desc, ok := iface.(*Descriptor); ok {
		return desc, true
	}

	// Check if it's a Descriptor value
	if desc, ok := iface.(Descriptor); ok {
		return &desc, true
	}

	// Try to manually extract fields if type names differ but structure matches
	if v.Kind() == reflect.Struct {
		desc := &Descriptor{}

		nameField := v.FieldByName("Name")
		if nameField.IsValid() && nameField.Kind() == reflect.String {
			desc.Name = nameField.String()
		}

		descField := v.FieldByName("Description")
		if descField.IsValid() && descField.Kind() == reflect.String {
			desc.Description = descField.String()
		}

		scenariosField := v.FieldByName("Scenarios")
		if scenariosField.IsValid() && scenariosField.Kind() == reflect.Slice {
			for i := 0; i < scenariosField.Len(); i++ {
				elem := scenariosField.Index(i)
				if scenarioDesc, ok := extractScenarioDescriptor(elem); ok {
					desc.Scenarios = append(desc.Scenarios, scenarioDesc)
				}
			}
		}

		if desc.Name != "" {
			return desc, true
		}
	}

	return nil, false
}

// extractScenarioDescriptor attempts to extract a scenario.Descriptor from a reflect.Value.
func extractScenarioDescriptor(v reflect.Value) (*scenario.Descriptor, bool) {
	if !v.IsValid() {
		return nil, false
	}

	// Handle pointer types
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	iface := v.Interface()

	// Check if it's already the correct type
	if desc, ok := iface.(*scenario.Descriptor); ok {
		return desc, true
	}

	if desc, ok := iface.(scenario.Descriptor); ok {
		return &desc, true
	}

	return nil, false
}

// memoryFS implements fs.FS using an in-memory file store.
type memoryFS struct {
	files map[string]*memoryFile
	dirs  map[string]bool
}

type memoryFile struct {
	name    string
	content []byte
	mode    fs.FileMode
	modTime time.Time
}

func newMemoryFS() *memoryFS {
	return &memoryFS{
		files: make(map[string]*memoryFile, 64),
		dirs:  make(map[string]bool, 16),
	}
}

func (m *memoryFS) addFile(path string, content []byte, mode fs.FileMode) {
	// Normalize path
	path = filepath.ToSlash(path)
	path = strings.TrimPrefix(path, "/")

	m.files[path] = &memoryFile{
		name:    filepath.Base(path),
		content: content,
		mode:    mode,
		modTime: time.Now(),
	}

	// Create parent directories
	dir := filepath.Dir(path)
	for dir != "." && dir != "" {
		m.dirs[dir] = true
		dir = filepath.Dir(dir)
	}
}

func (m *memoryFS) Open(name string) (fs.File, error) {
	// Normalize the path
	name = filepath.ToSlash(name)
	name = strings.TrimPrefix(name, "/")

	// Check if it's a file
	if f, ok := m.files[name]; ok {
		return &memoryFileReader{
			file:   f,
			reader: bytes.NewReader(f.content),
		}, nil
	}

	// Check if it's a directory
	if m.dirs[name] || name == "." || name == "" {
		return &memoryDirReader{
			fs:   m,
			name: name,
		}, nil
	}

	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
}

// memoryFileReader implements fs.File for reading file content.
type memoryFileReader struct {
	file   *memoryFile
	reader *bytes.Reader
}

func (f *memoryFileReader) Stat() (fs.FileInfo, error) {
	return &memoryFileInfo{file: f.file}, nil
}

func (f *memoryFileReader) Read(b []byte) (int, error) {
	return f.reader.Read(b)
}

func (f *memoryFileReader) Close() error {
	return nil
}

// memoryFileInfo implements fs.FileInfo.
type memoryFileInfo struct {
	file *memoryFile
}

func (fi *memoryFileInfo) Name() string       { return fi.file.name }
func (fi *memoryFileInfo) Size() int64        { return int64(len(fi.file.content)) }
func (fi *memoryFileInfo) Mode() fs.FileMode  { return fi.file.mode }
func (fi *memoryFileInfo) ModTime() time.Time { return fi.file.modTime }
func (fi *memoryFileInfo) IsDir() bool        { return false }
func (fi *memoryFileInfo) Sys() any           { return nil }

// memoryDirReader implements fs.File for directories.
type memoryDirReader struct {
	fs   *memoryFS
	name string
}

func (d *memoryDirReader) Stat() (fs.FileInfo, error) {
	return &memoryDirInfo{name: filepath.Base(d.name)}, nil
}

func (d *memoryDirReader) Read([]byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: d.name, Err: fs.ErrInvalid}
}

func (d *memoryDirReader) Close() error {
	return nil
}

// memoryDirInfo implements fs.FileInfo for directories.
type memoryDirInfo struct {
	name string
}

func (di *memoryDirInfo) Name() string       { return di.name }
func (di *memoryDirInfo) Size() int64        { return 0 }
func (di *memoryDirInfo) Mode() fs.FileMode  { return fs.ModeDir | 0755 }
func (di *memoryDirInfo) ModTime() time.Time { return time.Now() }
func (di *memoryDirInfo) IsDir() bool        { return true }
func (di *memoryDirInfo) Sys() any           { return nil }

// Ensure memoryFS implements fs.ReadDirFS for yaegi directory listing.
var _ fs.ReadDirFS = (*memoryFS)(nil)

func (m *memoryFS) ReadDir(name string) ([]fs.DirEntry, error) {
	name = filepath.ToSlash(name)
	name = strings.TrimPrefix(name, "/")

	if name == "" {
		name = "."
	}

	// Collect entries in this directory
	entries := make(map[string]fs.DirEntry, 16)

	prefix := name
	if prefix != "." && prefix != "" {
		prefix += "/"
	} else {
		prefix = ""
	}

	// Add files
	for path, file := range m.files {
		if !strings.HasPrefix(path, prefix) {
			continue
		}

		relPath := strings.TrimPrefix(path, prefix)
		parts := strings.SplitN(relPath, "/", 2)

		if len(parts) == 1 {
			// Direct file in this directory
			entries[parts[0]] = &memoryDirEntry{
				name:  file.name,
				isDir: false,
				mode:  file.mode,
			}
		} else {
			// Subdirectory
			entries[parts[0]] = &memoryDirEntry{
				name:  parts[0],
				isDir: true,
				mode:  fs.ModeDir | 0755,
			}
		}
	}

	// Add explicit directories
	for dirPath := range m.dirs {
		if !strings.HasPrefix(dirPath, prefix) {
			continue
		}

		relPath := strings.TrimPrefix(dirPath, prefix)
		parts := strings.SplitN(relPath, "/", 2)

		if _, exists := entries[parts[0]]; !exists {
			entries[parts[0]] = &memoryDirEntry{
				name:  parts[0],
				isDir: true,
				mode:  fs.ModeDir | 0755,
			}
		}
	}

	// Convert map to slice
	result := make([]fs.DirEntry, 0, len(entries))
	for _, entry := range entries {
		result = append(result, entry)
	}

	return result, nil
}

// memoryDirEntry implements fs.DirEntry.
type memoryDirEntry struct {
	name  string
	isDir bool
	mode  fs.FileMode
}

func (e *memoryDirEntry) Name() string               { return e.name }
func (e *memoryDirEntry) IsDir() bool                { return e.isDir }
func (e *memoryDirEntry) Type() fs.FileMode          { return e.mode.Type() }
func (e *memoryDirEntry) Info() (fs.FileInfo, error) { return e, nil }

// Implement fs.FileInfo for memoryDirEntry
func (e *memoryDirEntry) Size() int64        { return 0 }
func (e *memoryDirEntry) Mode() fs.FileMode  { return e.mode }
func (e *memoryDirEntry) ModTime() time.Time { return time.Now() }
func (e *memoryDirEntry) Sys() any           { return nil }
