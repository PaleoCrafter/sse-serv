// MIT license, dtg [at] lengo [dot] org · 10/2018

package amqp

import (
	"net/url"

	impl "github.com/streadway/amqp"
)

// Broker ...
type Broker interface {
	Connect() error
	Disconnect() error
	OpenChannel() (*impl.Channel, error)
	CloseChannel(c *impl.Channel) error
}

type broker struct {
	address    *url.URL
	connection *impl.Connection
}

func newBroker(url *url.URL) Broker {
	return &broker{address: url}
}

func (b *broker) Connect() error {
	var conn *impl.Connection
	var err error

	if conn, err = impl.Dial(b.address.String()); err == nil {
		b.connection = conn
	}

	return err
}

func (b *broker) Disconnect() error {
	return b.connection.Close()
}

func (b *broker) OpenChannel() (*impl.Channel, error) {
	return b.connection.Channel()
}

func (b *broker) CloseChannel(c *impl.Channel) error {
	return c.Close()
}
