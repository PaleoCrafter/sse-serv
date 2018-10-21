// MIT license, dtg [at] lengo [dot] org · 10/2018

package main

import (
	"flag"
	"fmt"
)

// GetOpt ...
type GetOpt interface {
	Complement(s State) bool
}

type getOpt struct {
	version *bool
}

// NewGetOpt ...
func NewGetOpt() GetOpt {
	options := &getOpt{
		version: flag.Bool("v", false, "display version"),
	}
	flag.Parse()
	return options
}

func (o *getOpt) Complement(s State) bool {
	if *o.version {
		fmt.Printf("%s\n", s.version.Semantic())
		return false
	}
	return true
}
