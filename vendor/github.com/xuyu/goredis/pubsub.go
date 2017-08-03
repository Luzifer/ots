package goredis

import (
	"errors"
	"strconv"
	"strings"
)

// Publish posts a message to the given channel.
// Integer reply: the number of clients that received the message.
func (r *Redis) Publish(channel, message string) (int64, error) {
	rp, err := r.ExecuteCommand("PUBLISH", channel, message)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// PubSub doc: http://redis.io/topics/pubsub
type PubSub struct {
	redis *Redis
	conn  *connection

	Patterns map[string]bool
	Channels map[string]bool
}

// PubSub new a PubSub from *redis.
func (r *Redis) PubSub() (*PubSub, error) {
	c, err := r.pool.Get()
	if err != nil {
		return nil, err
	}
	return &PubSub{
		redis:    r,
		conn:     c,
		Patterns: make(map[string]bool),
		Channels: make(map[string]bool),
	}, nil
}

// Close closes current pubsub command.
func (p *PubSub) Close() error {
	return p.conn.Conn.Close()
}

// Receive returns the reply of pubsub command.
// A message is a Multi-bulk reply with three elements.
// The first element is the kind of message:
// 1) subscribe: means that we successfully subscribed to the channel given as the second element in the reply.
// The third argument represents the number of channels we are currently subscribed to.
// 2) unsubscribe: means that we successfully unsubscribed from the channel given as second element in the reply.
// third argument represents the number of channels we are currently subscribed to.
// When the last argument is zero, we are no longer subscribed to any channel,
// and the client can issue any kind of Redis command as we are outside the Pub/Sub state.
// 3) message: it is a message received as result of a PUBLISH command issued by another client.
// The second element is the name of the originating channel, and the third argument is the actual message payload.
func (p *PubSub) Receive() ([]string, error) {
	rp, err := p.conn.RecvReply()
	if err != nil {
		return nil, err
	}
	command, err := rp.Multi[0].StringValue()
	if err != nil {
		return nil, err
	}
	switch strings.ToLower(command) {
	case "psubscribe", "punsubscribe":
		pattern, err := rp.Multi[1].StringValue()
		if err != nil {
			return nil, err
		}
		if command == "psubscribe" {
			p.Patterns[pattern] = true
		} else {
			delete(p.Patterns, pattern)
		}
		number, err := rp.Multi[2].IntegerValue()
		if err != nil {
			return nil, err
		}
		return []string{command, pattern, strconv.FormatInt(number, 10)}, nil
	case "subscribe", "unsubscribe":
		channel, err := rp.Multi[1].StringValue()
		if err != nil {
			return nil, err
		}
		if command == "subscribe" {
			p.Channels[channel] = true
		} else {
			delete(p.Channels, channel)
		}
		number, err := rp.Multi[2].IntegerValue()
		if err != nil {
			return nil, err
		}
		return []string{command, channel, strconv.FormatInt(number, 10)}, nil
	case "pmessage":
		pattern, err := rp.Multi[1].StringValue()
		if err != nil {
			return nil, err
		}
		channel, err := rp.Multi[2].StringValue()
		if err != nil {
			return nil, err
		}
		message, err := rp.Multi[3].StringValue()
		if err != nil {
			return nil, err
		}
		return []string{command, pattern, channel, message}, nil
	case "message":
		channel, err := rp.Multi[1].StringValue()
		if err != nil {
			return nil, err
		}
		message, err := rp.Multi[2].StringValue()
		if err != nil {
			return nil, err
		}
		return []string{command, channel, message}, nil
	}
	return nil, errors.New("pubsub protocol error")
}

// Subscribe channel [channel ...]
func (p *PubSub) Subscribe(channels ...string) error {
	args := packArgs("SUBSCRIBE", channels)
	return p.conn.SendCommand(args...)
}

// PSubscribe pattern [pattern ...]
func (p *PubSub) PSubscribe(patterns ...string) error {
	args := packArgs("PSUBSCRIBE", patterns)
	return p.conn.SendCommand(args...)
}

// UnSubscribe [channel [channel ...]]
func (p *PubSub) UnSubscribe(channels ...string) error {
	args := packArgs("UNSUBSCRIBE", channels)
	return p.conn.SendCommand(args...)
}

// PUnSubscribe [pattern [pattern ...]]
func (p *PubSub) PUnSubscribe(patterns ...string) error {
	args := packArgs("PUNSUBSCRIBE", patterns)
	return p.conn.SendCommand(args...)
}
