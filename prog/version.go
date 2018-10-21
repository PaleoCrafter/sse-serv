// MIT license, dtg [at] lengo [dot] org · 10/2018

package prog

import (
	"fmt"
)

// Version ...
type Version interface {
	Semantic() string
	Full() string
}

type version struct {
	nomen string
	major int
	minor int
	patch int
}

// NewVersion ...
func NewVersion(nomen string, major int, minor int, patch int) Version {
	return &version{
		nomen: nomen,
		major: major,
		minor: minor,
		patch: patch,
	}
}

func (v *version) Semantic() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}

func (v *version) Full() string {
	return fmt.Sprintf("%s@%s", v.nomen, v.Semantic())
}
