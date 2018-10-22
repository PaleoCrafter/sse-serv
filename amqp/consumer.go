// MIT license, dtg [at] lengo [dot] org · 10/2018

package amqp

import (
	"errors"

	impl "github.com/streadway/amqp"
)

// Provider ...
type Provider interface {
	Consumer(onCancel chan string, onClose chan *impl.Error) (Consumer, error)
}

// Consumer ...
type Consumer interface {
	Consume(queue string) (<-chan impl.Delivery, error)
	Close()
}

type provider struct {
	pool BrokerPool
}

type consumer struct {
	broker   Broker
	channel  *impl.Channel
	onCancel chan string
}

// NewProvider ...
func NewProvider(pool BrokerPool) Provider {
	return &provider{pool}
}

func (p *provider) Consumer(onCancel chan string, onClose chan *impl.Error) (Consumer, error) {
	var channel *impl.Channel
	var broker Broker
	var err error

	if broker, err = p.pool.Fetch(); err != nil {
		return nil, err
	}

	if channel, err = broker.OpenChannel(); err == nil {
		channel.NotifyClose(onClose)
	} else {
		return nil, err
	}

	return newConsumer(broker, channel, onCancel), nil
}

func newConsumer(broker Broker, channel *impl.Channel, onCancel chan string) Consumer {
	return &consumer{broker, channel, onCancel}
}

func (c *consumer) Consume(queue string) (<-chan impl.Delivery, error) {
	var t, f = true, false
	var que impl.Queue
	var err error

	if que, err = c.channel.QueueDeclare(queue, t, f, f, f, nil); err != nil {
		return nil, err
	}

	if que.Consumers == 0 {
		c.channel.NotifyCancel(c.onCancel)
		return c.channel.Consume(queue, "", f, f, f, f, nil)
	}

	return nil, errors.New("max consumers exceeded")
}

func (c *consumer) Close() {
	c.broker.CloseChannel(c.channel)
}
