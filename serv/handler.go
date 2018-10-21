// MIT license, dtg [at] lengo [dot] org · 10/2018

package serv

import (
	"fmt"
	"net/http"
	"time"

	impl "github.com/streadway/amqp"
	"sse-serv/amqp"
	"sse-serv/logg"
)

// ResponseHandler ...
type ResponseHandler interface {
	Handle()
}

type contextDone = func() <-chan struct{}

type responseHandler struct {
	writer   http.ResponseWriter
	request  http.Request
	provider amqp.Provider
	logger   logg.Logger
	created  time.Time
	pattern  QueuePattern
	creator  EventCreator
	ctxDone  contextDone
	headers  map[string]string
	queue    string
	events   int
}

// NewResponseHandler ...
func NewResponseHandler(
	a http.ResponseWriter,
	b http.Request,
	c amqp.Provider,
	d logg.Logger,
	e QueuePattern,
	f EventCreator,
	g contextDone,
	h map[string]string,
) ResponseHandler {

	return &responseHandler{
		writer:   a,
		request:  b,
		provider: c,
		logger:   d,
		pattern:  e,
		creator:  f,
		ctxDone:  g,
		headers:  h,
		created:  time.Now(),
	}
}

func (r *responseHandler) Handle() {
	var err error

	if r.request.Method != "GET" {
		r.sendStatus(405)
		return
	}
	if r.queue, err = r.pattern.QueueName(&r.request); err != nil {
		r.sendStatus(503)
		return
	}
	if err = r.drain(); err != nil {
		r.sendStatus(503)
	}
}

func (r *responseHandler) drain() error {
	var consumer amqp.Consumer
	var msgs <-chan impl.Delivery
	var err error

	onAmqpClose := make(chan *impl.Error)

	if consumer, err = r.provider.Consumer(onAmqpClose); err == nil {
		defer consumer.Close()
	} else {
		return err
	}

	if msgs, err = consumer.Consume(r.queue); err != nil {
		return err
	}

	r.setHeaders()
	r.sendBanner()

	go func() {
		for delivery := range msgs {
			event := r.creator.CreateEvent(delivery.Body)
			r.sendEvent(event)
		}
	}()

	r.logger.Info("%s %s", "connected", r.queue)

	select {
	case _ = <-onAmqpClose:
	case _ = <-r.ctxDone():
	}

	d := time.Now().Sub(r.created).Truncate(1e7).String()
	r.logger.Info("disconnected (%s, %d events)", d, r.events)

	return nil
}

func (r *responseHandler) sendStatus(status int) {
	r.writer.WriteHeader(status)
}

func (r *responseHandler) sendBanner() {
	fmt.Fprintf(r.writer, ": SSE stream\n")
	r.commit()
}

func (r *responseHandler) sendEvent(ev Event) {
	fmt.Fprint(r.writer, ev.String())
	r.commit()
	r.events++
}

func (r *responseHandler) sendPulse() {
	fmt.Fprintf(r.writer, ":\n")
	r.commit()
}

func (r *responseHandler) commit() {
	r.writer.Write([]byte("\n"))
	r.writer.(http.Flusher).Flush()
}

func (r *responseHandler) setHeaders() {
	for k, v := range r.headers {
		r.writer.Header().Set(k, v)
	}
}
