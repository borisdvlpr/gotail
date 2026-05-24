//go:build !linux && !darwin && !windows

package file

import "github.com/spf13/afero"

func NewSystemSearcher(fsys afero.Fs) *SystemSearcher {
	return &SystemSearcher{
		Fsys:     fsys,
		Searcher: nil, // sentinel: unsupported OS
	}
}
