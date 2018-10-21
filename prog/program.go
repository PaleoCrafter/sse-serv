// MIT license, dtg [at] lengo [dot] org · 10/2018

package prog

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"sse-serv/logg"
	"sse-serv/serv"
)

// Program ...
type Program interface {
	Run() error
}

type program struct {
	server serv.Server
	logger logg.Logger
}

// NewProgram ...
func NewProgram(server serv.Server, logger logg.Logger) Program {
	return &program{
		server: server,
		logger: logger,
	}
}

func (p *program) Run() error {

	termSignal := make(chan bool, 1)
	signalName := make(chan os.Signal, 1)
	signal.Notify(signalName, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		name := <-signalName
		p.logger.Info("caught signal '%s', bailing out", name)
		termSignal <- true
	}()

	haveServed := make(chan bool, 1)

	go func() {
		haveServed <- p.server.Start()
	}()

	select {
	case _ = <-termSignal:
		// FIXME shutdown server

	case served := <-haveServed:
		if !served {
			return errors.New("the server made a booboo")
		}
	}

	return nil
}
