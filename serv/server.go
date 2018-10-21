// MIT license, dtg [at] lengo [dot] org · 10/2018

package serv

import (
	"net/http"

	"sse-serv/logg"
)

// Server ...
type Server interface {
	Start() bool
	Stop() bool
}

type server struct {
	address string
	handler http.HandlerFunc
	logger  logg.Logger
	server  http.Server
}

// NewServer ...
func NewServer(
	address string,
	handler http.HandlerFunc,
	logger logg.Logger,
) Server {

	return &server{
		address: address,
		handler: handler,
		logger:  logger,
		server:  http.Server{Addr: address, Handler: handler},
	}
}

func (s *server) Start() bool {
	var err error

	s.logger.Info("listening %s", s.address)

	if err = s.server.ListenAndServe(); err == http.ErrServerClosed {
		s.logger.Info("shutdown %s", s.address)
	} else {
		s.logger.Error("%s", err)
	}

	return err == http.ErrServerClosed
}

func (s *server) Stop() bool {
	return false
}
