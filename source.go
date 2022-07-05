package nrelay

import (
	"log"

	"github.com/nats-io/nats.go"
	"github.com/octu0/chanque"
	"github.com/pkg/errors"
)

type Source interface {
	Open() error
	Close() error
	Subscribe(topic string, prefixSize int, workers []chanque.Worker) error
	Unsubscribe() error
}

// check interface
var (
	_ (Source) = (*MultipleSource)(nil)
)

type MultipleSource struct {
	natsUrls []string
	natsOpts []nats.Option
	logger   *log.Logger
	conns    []*nats.Conn
	subs     []*nats.Subscription
	opts     *sourceOptions
}

func (s *MultipleSource) Open() error {
	conns := make([]*nats.Conn, len(s.natsUrls))
	for i, url := range s.natsUrls {
		conn, err := nats.Connect(url, s.natsOpts...)
		if err != nil {
			return errors.WithStack(err)
		}
		s.logger.Printf("debug: source connect %s", url)

		conns[i] = conn
	}
	s.conns = conns
	return nil
}

func (s *MultipleSource) Subscribe(topic string, prefixSize int, workers []chanque.Worker) error {
	dist := newDistribute(workers)
	subs := make([]*nats.Subscription, len(s.conns))
	for i, conn := range s.conns {
		sub, err := conn.Subscribe(topic, s.createSubscribeHandler(prefixSize, dist))
		if err != nil {
			return errors.WithStack(err)
		}
		subs[i] = sub
	}
	for _, conn := range s.conns {
		conn.Flush()
	}
	s.subs = subs
	return nil
}

func (s *MultipleSource) Unsubscribe() error {
	for _, sub := range s.subs {
		if err := sub.Unsubscribe(); err != nil {
			return errors.WithStack(err)
		}
	}
	s.subs = s.subs[len(s.subs):]
	return nil
}

func (s *MultipleSource) Close() error {
	if err := s.Unsubscribe(); err != nil {
		return errors.WithStack(err)
	}

	for _, conn := range s.conns {
		conn.Close()
	}
	return nil
}

func (s *MultipleSource) createSubscribeHandler(prefixSize int, dist *distribute) nats.MsgHandler {
	return func(msg *nats.Msg) {
		s.opts.OnMessageReceived(msg)

		if 0 < prefixSize {
			if ok := dist.Publish(msg.Subject[0:prefixSize], msg); ok != true {
				s.logger.Printf("warn: failed to publish: %s", msg.Subject)
				s.opts.OnMessageRelayFailed(msg)
				return
			}

			s.opts.OnMessageRelayed(msg)
			return
		}

		if ok := dist.Publish(msg.Subject, msg); ok != true {
			s.logger.Printf("warn: failed to publish: %s", msg.Subject)
			s.opts.OnMessageRelayFailed(msg)
			return
		}

		s.opts.OnMessageRelayed(msg)
	}
}

func NewMultipleSource(urls []string, sourceOpts []SourceOption, natsOpts []nats.Option, logger *log.Logger) *MultipleSource {
	opts := &sourceOptions{}
	for _, o := range sourceOpts {
		o(opts)
	}

	return &MultipleSource{urls, natsOpts, logger, nil, nil, opts}
}
