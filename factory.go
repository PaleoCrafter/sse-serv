// MIT license, dtg [at] lengo [dot] org · 10/2018

package main

import (
	"net/http"

	"sse-serv/amqp"
	"sse-serv/logg"
	"sse-serv/prog"
	"sse-serv/serv"
)

// Factory ...
type Factory interface {
	CreateProgram() prog.Program
}

type factory struct {
	state      State
	brokerPool amqp.BrokerPool
	healthLog  logg.Logger
	accessLog  logg.Logger
}

// NewFactory ...
func NewFactory(s State) Factory {
	return &factory{state: s}
}

func (f *factory) Config() Config {
	return f.state.config
}

func (f *factory) CreateProgram() prog.Program {
	return prog.NewProgram(
		f.createEventServer(),
		f.createHealthLogger(),
	)
}

func (f *factory) createEventServer() serv.Server {
	return serv.NewServer(
		f.Config().Server.Listen,
		f.createResponseHandler(),
		f.createHealthLogger(),
	)
}

func (f *factory) createResponseHandler() http.HandlerFunc {
	return func() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			waitForContextDone :=
				func() <-chan struct{} {
					return r.Context().Done()
				}

			h := serv.NewResponseHandler(
				w,
				*r,
				f.createConsumerProvider(),
				f.createAccessLogger(),
				f.createQueuePattern(),
				f.createEventCreator(),
				waitForContextDone,
				f.Config().Header.SSE,
			)

			h.Handle()
		}
	}()
}

func (f *factory) createBrokerPool() amqp.BrokerPool {
	if f.brokerPool == nil {
		f.brokerPool = amqp.NewBrokerPool(
			f.Config().Broker.Connect,
			f.createHealthLogger(),
		)
		f.brokerPool.WarmUp()
	}
	return f.brokerPool
}

func (f *factory) createQueuePattern() serv.QueuePattern {
	return serv.NewPattern(
		f.Config().Queue.Pattern,
	)
}

func (f *factory) createEventCreator() serv.EventCreator {
	return serv.NewEventCreator()
}

func (f *factory) createConsumerProvider() amqp.Provider {
	return amqp.NewProvider(
		f.createBrokerPool(),
	)
}

func (f *factory) createHealthLogger() logg.Logger {
	if f.healthLog == nil {
		f.healthLog = f.createLogger(
			f.Config().Logger.Health.File,
		)
	}
	return f.healthLog
}

func (f *factory) createAccessLogger() logg.Logger {
	if f.accessLog == nil {
		f.accessLog = f.createLogger(
			f.Config().Logger.Access.File,
		)
	}
	return f.accessLog
}

func (f *factory) createLogger(filename string) logg.Logger {
	logger := logg.NewLogger(filename)
	if err := logger.Open(); err != nil {
		panic(err)
	} else {
		return logger
	}
}
