// Package plugin provides dynamic plugin loading using Yaegi.
package plugin

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"
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

// PluginLoader handles loading plugins from various sources at runtime.
type PluginLoader struct {
	logger           *logrus.Entry
	pluginRegistry   *PluginRegistry
	scenarioRegistry *ScenarioRegistry
	mu               sync.Mutex
	cleanupFn        func(*LoadedPlugin) // callback for cleanup notifications
}

// NewPluginLoader creates a new plugin loader with the given registries.
func NewPluginLoader(
	logger logrus.FieldLogger,
	pluginRegistry *PluginRegistry,
	scenarioRegistry *ScenarioRegistry,
) *PluginLoader {
	return &PluginLoader{
		logger:           logger.WithField("component", "plugin-loader"),
		pluginRegistry:   pluginRegistry,
		scenarioRegistry: scenarioRegistry,
	}
}

// SetCleanupCallback sets a callback function to be called when a plugin
// may be ready for cleanup (e.g., when a scenario is replaced or unregistered).
func (l *PluginLoader) SetCleanupCallback(fn func(*LoadedPlugin)) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.cleanupFn = fn
}

// LoadFromBytes loads a plugin from tar(.gz) bytes.
// The plugin name is determined from plugin.yaml inside the archive.
// If compressed is true, the data is treated as gzip-compressed.
func (l *PluginLoader) LoadFromBytes(data []byte, compressed bool) (*LoadedPlugin, error) {
	return l.LoadFromReader(bytes.NewReader(data), compressed)
}

// LoadFromFile loads a plugin from a tar(.gz) file path.
// The plugin name is determined from plugin.yaml inside the archive.
// Compression is auto-detected based on the .gz extension.
func (l *PluginLoader) LoadFromFile(filePath string) (*LoadedPlugin, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin file: %w", err)
	}
	defer file.Close()

	// Determine compression from filename
	baseName := filepath.Base(filePath)
	compressed := strings.HasSuffix(baseName, ".gz")

	// Plugin name will be determined from plugin.yaml
	loaded, err := l.LoadFromReader(file, compressed)
	if err != nil {
		return nil, err
	}

	loaded.SourceType = PluginSourceFile

	return loaded, nil
}

// LoadFromReader loads a plugin from a tar(.gz) stream.
// It first extracts and parses plugin.yaml to get the plugin name,
// then extracts the full archive to the appropriate directory structure.
// If compressed is true, the stream is treated as gzip-compressed.
func (l *PluginLoader) LoadFromReader(data io.Reader, compressed bool) (*LoadedPlugin, error) {
	// Read all data into memory so we can scan for plugin.yaml first
	allData, err := io.ReadAll(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin data: %w", err)
	}

	// Extract and parse plugin.yaml from the archive
	metadata, err := l.extractMetadataFromTar(allData, compressed)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin metadata: %w", err)
	}

	// Use name from metadata instead of filename
	actualPluginName := metadata.Name

	l.logger.Debugf("loading plugin '%s' (build: %s, version: %s)",
		metadata.Name, metadata.BuildTime, metadata.GitVersion)

	// Now extract the full archive
	var tarReader io.Reader = bytes.NewReader(allData)

	if compressed {
		gzReader, err := gzip.NewReader(bytes.NewReader(allData))
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzReader.Close()

		tarReader = gzReader
	}

	// Create temp directory for this plugin using the actual plugin name
	tempDir, pluginPath, err := l.createTempDir(actualPluginName)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Extract tar to the plugin path
	err = l.extractTar(tarReader, pluginPath)
	if err != nil {
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("failed to extract plugin: %w", err)
	}

	// Load from the temp directory filesystem
	filesys := NewSymlinkFS(tempDir)
	desc, err := l.loadFromFS(actualPluginName, filesys)
	if err != nil {
		os.RemoveAll(tempDir)
		return nil, err
	}

	loaded := NewLoadedPlugin(desc, metadata, tempDir, pluginPath, PluginSourceBytes)

	return loaded, nil
}

// extractMetadataFromTar extracts and parses plugin.yaml from a tar archive.
func (l *PluginLoader) extractMetadataFromTar(data []byte, compressed bool) (*PluginMetadata, error) {
	var tarReader io.Reader = bytes.NewReader(data)

	if compressed {
		gzReader, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzReader.Close()

		tarReader = gzReader
	}

	tr := tar.NewReader(tarReader)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("failed to read tar entry: %w", err)
		}

		// Look for plugin.yaml at the root level
		cleanPath := filepath.Clean(header.Name)
		cleanPath = strings.TrimPrefix(cleanPath, "./")
		cleanPath = strings.TrimPrefix(cleanPath, "/")

		if cleanPath == PluginMetadataFile && header.Typeflag == tar.TypeReg {
			content, err := io.ReadAll(tr)
			if err != nil {
				return nil, fmt.Errorf("failed to read %s: %w", PluginMetadataFile, err)
			}

			return ParsePluginMetadata(content)
		}
	}

	return nil, fmt.Errorf("%s not found in plugin archive (required for tar.gz plugins)", PluginMetadataFile)
}

// LoadFromURL loads a plugin from a remote URL (tar.gz).
// The plugin name is determined from plugin.yaml inside the archive.
func (l *PluginLoader) LoadFromURL(url string) (*LoadedPlugin, error) {
	l.logger.Infof("downloading plugin from URL: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download plugin: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download plugin: HTTP %d", resp.StatusCode)
	}

	// Plugin name will be determined from plugin.yaml
	loaded, err := l.LoadFromReader(resp.Body, true)
	if err != nil {
		return nil, err
	}

	loaded.SourceType = PluginSourceURL

	return loaded, nil
}

// LoadFromLocalPath loads a plugin from a local directory path.
// This creates a symlink in the temp directory to the local path.
// Local path plugins use the directory name as plugin name and spamoor's
// build info for version metadata (no plugin.yaml required).
func (l *PluginLoader) LoadFromLocalPath(localPath string) (*LoadedPlugin, error) {
	// Resolve to absolute path
	absPath, err := filepath.Abs(localPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	// Check if path exists and is a directory
	info, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat path: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", absPath)
	}

	// Derive plugin name from directory name
	pluginName := filepath.Base(absPath)

	// Create metadata using directory name and spamoor's build info
	metadata := NewLocalPluginMetadata(pluginName)

	l.logger.Debugf("loading local plugin '%s' (using spamoor build: %s, version: %s)",
		metadata.Name, metadata.BuildTime, metadata.GitVersion)

	// Create temp directory
	tempDir, pluginPath, err := l.createTempDir(pluginName)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Create symlink from pluginPath to absPath
	err = os.Symlink(absPath, pluginPath)
	if err != nil {
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("failed to create symlink: %w", err)
	}

	// Load using SymlinkFS which follows symlinks
	filesys := NewSymlinkFS(tempDir)
	desc, err := l.loadFromFS(pluginName, filesys)
	if err != nil {
		os.RemoveAll(tempDir)
		return nil, err
	}

	loaded := NewLoadedPlugin(desc, metadata, tempDir, pluginPath, PluginSourceLocal)

	return loaded, nil
}

// createTempDir creates the temp directory structure for a plugin.
// Returns the base temp dir and the plugin path within it.
func (l *PluginLoader) createTempDir(pluginName string) (tempDir, pluginPath string, err error) {
	// Create base temp directory
	tempDir, err = os.MkdirTemp("", fmt.Sprintf("spamoor-plugin-%s-", pluginName))
	if err != nil {
		return "", "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Create the GOPATH structure: src/github.com/ethpandaops/spamoor/plugins/<plugin-name>
	pluginPath = filepath.Join(tempDir, PluginBasePath, pluginName)

	// Create parent directories (but not the pluginPath itself - it will be created/symlinked)
	parentPath := filepath.Dir(pluginPath)
	err = os.MkdirAll(parentPath, 0755)
	if err != nil {
		os.RemoveAll(tempDir)
		return "", "", fmt.Errorf("failed to create plugin directory structure: %w", err)
	}

	return tempDir, pluginPath, nil
}

// extractTar extracts a tar archive to the given destination directory.
func (l *PluginLoader) extractTar(tarReader io.Reader, destDir string) error {
	// Create destination directory
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	tr := tar.NewReader(tarReader)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("failed to read tar entry: %w", err)
		}

		// Clean the path and remove leading "./" or "/"
		cleanPath := filepath.Clean(header.Name)
		cleanPath = strings.TrimPrefix(cleanPath, "./")
		cleanPath = strings.TrimPrefix(cleanPath, "/")

		targetPath := filepath.Join(destDir, cleanPath)

		// Ensure the target path is within destDir (prevent path traversal)
		if !strings.HasPrefix(targetPath, destDir) {
			return fmt.Errorf("invalid tar path: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(targetPath, header.FileInfo().Mode())
			if err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
		case tar.TypeReg:
			// Create parent directories
			err = os.MkdirAll(filepath.Dir(targetPath), 0755)
			if err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}

			// Create file
			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, header.FileInfo().Mode())
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", targetPath, err)
			}

			_, err = io.Copy(outFile, tr)
			outFile.Close()
			if err != nil {
				return fmt.Errorf("failed to write file %s: %w", targetPath, err)
			}
		}
	}

	return nil
}

// loadFromFS interprets a plugin from a filesystem.
func (l *PluginLoader) loadFromFS(pluginName string, filesys fs.FS) (*scenario.PluginDescriptor, error) {
	// Create interpreter with the filesystem as GOPATH source
	i := l.newInterpreter(filesys)

	// Import the main plugin package
	pluginPkg := fmt.Sprintf("github.com/ethpandaops/spamoor/plugins/%s", pluginName)

	// First, import the package (using safeEval to catch panics)
	_, err := l.safeEval(i, fmt.Sprintf(`import plugin "%s"`, pluginPkg))
	if err != nil {
		return nil, l.wrapEvalError(pluginName, err)
	}

	// Then evaluate the PluginDescriptor (using safeEval to catch panics)
	v, err := l.safeEval(i, "plugin.PluginDescriptor")
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

// RegisterPluginScenarios registers all scenarios from a loaded plugin.
func (l *PluginLoader) RegisterPluginScenarios(loaded *LoadedPlugin) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Register plugin in the plugin registry
	oldPlugin := l.pluginRegistry.Register(loaded)

	// Register each scenario
	for _, scenarioDesc := range loaded.Descriptor.Scenarios {
		entry := &ScenarioEntry{
			Descriptor: scenarioDesc,
			Source:     ScenarioSourcePlugin,
			Plugin:     loaded,
		}

		oldEntry, err := l.scenarioRegistry.Register(entry)
		if err != nil {
			return fmt.Errorf("failed to register scenario %s: %w", scenarioDesc.Name, err)
		}

		// Track scenario in the loaded plugin
		loaded.AddScenario(scenarioDesc.Name)

		// If we replaced a scenario from a different plugin, update that plugin's tracking
		if oldEntry != nil && oldEntry.Plugin != nil && oldEntry.Plugin != loaded {
			oldEntry.Plugin.RemoveScenario(scenarioDesc.Name)
			l.maybeCleanup(oldEntry.Plugin)
		}

		l.logger.Infof("registered scenario from plugin: %s", scenarioDesc.Name)
	}

	// If we replaced an old plugin entirely, check if it can be cleaned up
	if oldPlugin != nil && oldPlugin != loaded {
		l.maybeCleanup(oldPlugin)
	}

	return nil
}

// UnregisterPluginScenarios removes all scenarios from a plugin.
func (l *PluginLoader) UnregisterPluginScenarios(pluginName string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	plugin := l.pluginRegistry.Get(pluginName)
	if plugin == nil {
		return fmt.Errorf("plugin not found: %s", pluginName)
	}

	// Remove each scenario from the registry
	for _, scenarioDesc := range plugin.Descriptor.Scenarios {
		_, err := l.scenarioRegistry.Remove(scenarioDesc.Name)
		if err != nil {
			l.logger.Warnf("failed to remove scenario %s: %v", scenarioDesc.Name, err)
		}

		plugin.RemoveScenario(scenarioDesc.Name)
	}

	// Remove plugin from registry
	l.pluginRegistry.Remove(pluginName)

	// Check if plugin can be cleaned up
	l.maybeCleanup(plugin)

	return nil
}

// maybeCleanup checks if a plugin can be cleaned up and triggers cleanup if so.
// Must be called with l.mu held.
func (l *PluginLoader) maybeCleanup(plugin *LoadedPlugin) {
	if plugin == nil || plugin.IsCleanedUp() {
		return
	}

	if plugin.CanCleanup() {
		if l.cleanupFn != nil {
			// Call cleanup callback asynchronously to avoid deadlock
			go l.cleanupFn(plugin)
		} else {
			// Default cleanup: remove temp directory
			go l.CleanupPlugin(plugin)
		}
	}
}

// CleanupPlugin removes the temp directory for a plugin.
func (l *PluginLoader) CleanupPlugin(plugin *LoadedPlugin) error {
	if plugin.IsCleanedUp() {
		return nil
	}

	plugin.MarkCleanedUp()

	if plugin.TempDir != "" {
		l.logger.Infof("cleaning up plugin temp directory: %s", plugin.TempDir)

		err := os.RemoveAll(plugin.TempDir)
		if err != nil {
			return fmt.Errorf("failed to remove temp directory: %w", err)
		}
	}

	return nil
}

// Shutdown cleans up all loaded plugins.
func (l *PluginLoader) Shutdown() {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, plugin := range l.pluginRegistry.GetAll() {
		if err := l.CleanupPlugin(plugin); err != nil {
			l.logger.Warnf("failed to cleanup plugin %s: %v", plugin.Descriptor.Name, err)
		}
	}
}

// GetPluginRegistry returns the plugin registry.
func (l *PluginLoader) GetPluginRegistry() *PluginRegistry {
	return l.pluginRegistry
}

// GetScenarioRegistry returns the scenario registry.
func (l *PluginLoader) GetScenarioRegistry() *ScenarioRegistry {
	return l.scenarioRegistry
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

// safeEval wraps interpreter Eval calls with panic recovery.
// Yaegi may panic on invalid plugins; this converts panics to errors.
func (l *PluginLoader) safeEval(i *interp.Interpreter, src string) (v reflect.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("interpreter panic: %v", r)
		}
	}()

	return i.Eval(src)
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
func extractPluginDescriptor(v reflect.Value) (*scenario.PluginDescriptor, bool) {
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
	if desc, ok := iface.(*scenario.PluginDescriptor); ok {
		return desc, true
	}

	// Check if it's a Descriptor value
	if desc, ok := iface.(scenario.PluginDescriptor); ok {
		return &desc, true
	}

	// Try to manually extract fields if type names differ but structure matches
	if v.Kind() == reflect.Struct {
		desc := &scenario.PluginDescriptor{}

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
