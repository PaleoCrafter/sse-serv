// MIT license, dtg [at] lengo [dot] org · 10/2018

package serv

import (
	"net/http"
	"testing"
)

// Create instance
func TestNewServer(t *testing.T) {
	if r := NewServer("", nil, nil); r == nil {
		t.Error("expected constructor to work")
	}
}

// Server must listen
func TestServer_Start(t *testing.T) {
	_ = &server{
		address: ":1234",
		handler: nil,
		logger:  nil,
		server:  http.Server{},
	}

	// 	s.Start()
}
