//go:build darwin

package file

import "github.com/spf13/afero"

// NewSystemSearcher returns a SystemSearcher for macOS. No DriveSearcher is
// needed here — FindUserData handles darwin by walking /Volumes directly, so
// Searcher is left nil as the sentinel for that path.
func NewSystemSearcher(fsys afero.Fs) *SystemSearcher {
	return &SystemSearcher{
		Fsys:     fsys,
		Searcher: nil,
	}
}
