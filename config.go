// MIT license, dtg [at] lengo [dot] org · 10/2018

package main

import (
	"encoding/json"
	"io/ioutil"
)

// Config ...
type Config struct {
	Server struct {
		Listen string
	}
	Broker struct {
		Connect []string
	}
	Queue struct {
		Pattern string
		Expires string
	}
	Logger struct {
		Health struct {
			File string
		}
		Access struct {
			File string
		}
	}
	Header struct {
		CORS map[string]string
		SSE  map[string]string
	}
}

// NewConfig ...
func NewConfig(filename string) Config {
	var (
		bytes []byte
		err   error
	)
	if bytes, err = ioutil.ReadFile(filename); err != nil {
		panic(err)
	}

	config := &Config{}
	if err = json.Unmarshal(bytes, config); err != nil {
		panic(err)
	}

	return *config
}
