// MIT license, dtg [at] lengo [dot] org · 10/2018

package amqp

import (
	"errors"
	"net/url"

	"sse-serv/logg"
)

// BrokerPool ...
type BrokerPool interface {
	WarmUp() error
	Fetch() (Broker, error)
}

type brokerPool struct {
	urls   []string
	items  []poolItem
	logger logg.Logger
}

type poolItem struct {
	url  string
	conn Broker
}

// NewBrokerPool ...
func NewBrokerPool(urls []string, logger logg.Logger) BrokerPool {
	return &brokerPool{urls: urls, logger: logger}
}

func (p *brokerPool) WarmUp() error {
	for _, u := range p.urls {

		uri, _ := url.Parse(u)
		conn := newBroker(uri)

		if err := conn.Connect(); err != nil {
			p.logger.Warn("%s", err)
			continue
		}

		p.items = append(
			p.items,
			poolItem{url: u, conn: conn},
		)

		safe := uri
		safe.User = nil

		p.logger.Info("%s %s", "connected", safe)
	}

	return nil
}

func (p *brokerPool) Fetch() (Broker, error) {
	if len(p.items) == 0 {
		return nil, errors.New("out of broker connections")
	}
	return p.items[0].conn, nil
}
