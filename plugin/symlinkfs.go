// Package plugin provides dynamic plugin loading using Yaegi.
package plugin

import (
	"io/fs"
	"os"
	"path/filepath"
)

// SymlinkFS implements fs.FS and fs.ReadDirFS while following symlinks.
// This is used for local plugin paths that may be symlinked.
type SymlinkFS struct {
	basePath string
}

// NewSymlinkFS creates a new SymlinkFS rooted at the given base path.
func NewSymlinkFS(basePath string) *SymlinkFS {
	return &SymlinkFS{
		basePath: basePath,
	}
}

// Open opens a file at the given path, following symlinks.
func (s *SymlinkFS) Open(name string) (fs.File, error) {
	// Validate the path
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}

	fullPath := filepath.Join(s.basePath, filepath.FromSlash(name))

	return os.Open(fullPath)
}

// ReadDir reads the directory at the given path, following symlinks.
func (s *SymlinkFS) ReadDir(name string) ([]fs.DirEntry, error) {
	// Validate the path
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: fs.ErrInvalid}
	}

	fullPath := filepath.Join(s.basePath, filepath.FromSlash(name))

	return os.ReadDir(fullPath)
}

// Stat returns file info for the given path, following symlinks.
func (s *SymlinkFS) Stat(name string) (fs.FileInfo, error) {
	// Validate the path
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "stat", Path: name, Err: fs.ErrInvalid}
	}

	fullPath := filepath.Join(s.basePath, filepath.FromSlash(name))

	return os.Stat(fullPath)
}

// Ensure SymlinkFS implements the required interfaces.
var (
	_ fs.FS        = (*SymlinkFS)(nil)
	_ fs.ReadDirFS = (*SymlinkFS)(nil)
	_ fs.StatFS    = (*SymlinkFS)(nil)
)
