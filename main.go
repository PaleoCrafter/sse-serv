// MIT license, dtg [at] lengo [dot] org · 10/2018

package main

import (
	"sse-serv/prog"
)

var nomen = "sse-serv"
var major = 0
var minor = 0
var patch = 0

// State ...
type State struct {
	version prog.Version
	config  Config
}

func main() {
	state := State{
		version: prog.NewVersion(nomen, major, minor, patch),
		config:  NewConfig("./config.json"),
	}

	if NewGetOpt().Complement(state) {
		p := NewFactory(state)

		if err := p.CreateProgram().Run(); err != nil {
			panic(err)
		}
	}
}
