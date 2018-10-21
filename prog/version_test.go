// MIT license, dtg [at] lengo [dot] org · 10/2018

package prog

import (
	"testing"
)

// Create instance
func TestNewVersion(t *testing.T) {
	if r := NewVersion("", 0, 0, 0); r == nil {
		t.Error("expected constructor to work")
	}
}

// Semantic version
func TestVersion_Semantic(t *testing.T) {
	r := NewVersion("foo", 1, 2, 3)
	if r.Semantic() != "1.2.3" {
		t.Error("expected semantic version pattern")
	}
}

// Full version with name
func TestVersion_Full(t *testing.T) {
	r := NewVersion("foo", 1, 2, 3)
	if r.Full() != "foo@1.2.3" {
		t.Error("expected full version pattern")
	}
}
