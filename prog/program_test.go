// MIT license, dtg [at] lengo [dot] org · 10/2018

package prog

import (
	"testing"
)

// Create instance
func TestNewProgram(t *testing.T) {
	if r := NewProgram(nil, nil); r == nil {
		t.Error("expected constructor to work")
	}
}
