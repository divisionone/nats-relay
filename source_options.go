package nrelay

import "github.com/nats-io/nats.go"

type SourceOption func(*sourceOptions)

type sourceOptions struct {
	// onMessageReceived is called when a subscription is received.
	onMessageReceived func(msg *nats.Msg)
	// onMessageRelayed is called when a message is relayed.
	onMessageRelayed func(msg *nats.Msg)
	// onMessageRelayFailed is called when a message is failed to relay.
	onMessageRelayFailed func(msg *nats.Msg)
}

// SourceOptOnSubscriptionReceived is a callback that is called when a subscription is received.
func SourceOptOnSubscriptionReceived(fn func(msg *nats.Msg)) SourceOption {
	return func(opt *sourceOptions) {
		opt.onMessageReceived = fn
	}
}

// SourceOptOnMessageRelayed is a callback that is called when a message is relayed.
func SourceOptOnMessageRelayed(fn func(msg *nats.Msg)) SourceOption {
	return func(opt *sourceOptions) {
		opt.onMessageRelayed = fn
	}
}

// SourceOptOnMessageRelayFailed is a callback that is called when a message is failed to relay.
func SourceOptOnMessageRelayFailed(fn func(msg *nats.Msg)) SourceOption {
	return func(opt *sourceOptions) {
		opt.onMessageRelayFailed = fn
	}
}

// OnMessageReceived is a callback that is called when a subscription is received.
func (o *sourceOptions) OnMessageReceived(msg *nats.Msg) {
	if o.onMessageReceived != nil {
		o.onMessageReceived(msg)
	}
}

// OnMessageRelayed is a callback that is called when a message is relayed.
func (o *sourceOptions) OnMessageRelayed(msg *nats.Msg) {
	if o.onMessageRelayed != nil {
		o.onMessageRelayed(msg)
	}
}

// OnMessageRelayFailed is a callback that is called when a message is failed to relay.
func (o *sourceOptions) OnMessageRelayFailed(msg *nats.Msg) {
	if o.onMessageRelayFailed != nil {
		o.onMessageRelayFailed(msg)
	}
}
