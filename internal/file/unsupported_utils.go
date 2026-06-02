//go:build !linux && !darwin && !windows

package file

import "github.com/spf13/afero"

// NewSystemSearcher returns a SystemSearcher for unrecognised operating systems.
// The embedded Searcher is set to nil as a sentinel value, indicating that drive
// searching is not supported on the current platform. Callers should check for
// this condition before invoking any search operations.
func NewSystemSearcher(fsys afero.Fs) *SystemSearcher {
	return &SystemSearcher{
		Fsys:     fsys,
		Searcher: nil, // sentinel: unsupported OS
	}
}
