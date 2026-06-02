//go:build darwin

package file

import (
	"context"

	"github.com/spf13/afero"
)

// DarwinSearcher implements DriveSearcher for macOS by walking /Volumes,
// where the OS automatically mounts removable media such as SD cards.
type DarwinSearcher struct{}

// NewSystemSearcher returns a SystemSearcher for macOS, wired with a
// DarwinSearcher that walks /Volumes to locate the target file.
func NewSystemSearcher(fsys afero.Fs) *SystemSearcher {
	return &SystemSearcher{
		Fsys:     fsys,
		Searcher: &DarwinSearcher{},
	}
}

// Search locates fileName by walking /Volumes, where macOS mounts removable
// drives. It returns the path of the first match, or an empty string if not
// found. No concurrency is needed as the search is a single synchronous walk.
func (ds *DarwinSearcher) Search(ctx context.Context, fsys afero.Fs, fileName string) (string, error) {
	return GetFilePath(fsys, "/Volumes", fileName)
}
