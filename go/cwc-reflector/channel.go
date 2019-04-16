package main

import (
	"net"
)
import (
	"../bitoip"
	"time"
)

type Subscriber struct {
	key bitoip.CarrierKeyType
	address net.Addr
	last_tx time.Time
}

type Channel struct {
	ChannelId bitoip.ChannelIdType
	Subscribers map[bitoip.CarrierKeyType]Subscriber
	Addresses map[string]Subscriber
	LastKey bitoip.CarrierKeyType;
}

var channels map[uint16]Channel = make(map[uint16]Channel)

func NewChannel(channel_id bitoip.ChannelIdType) Channel {
	return Channel {
		 channel_id,
		make(map[bitoip.CarrierKeyType]Subscriber),
		 make(map[string]Subscriber),
		   1,
	}
}

func ChannelIds() []uint16 {
	keys := make([]uint16, 0, len(channels))
	for k := range channels {
		keys = append(keys, k)
	}
	return keys
}
/**
 * get a channel by channel_id
 */
func GetChannel(channel_id bitoip.ChannelIdType) Channel {
	if channel, ok := channels[channel_id]; ok {
		return channel;
	} else {
		channels[channel_id] = NewChannel(channel_id)
		return channels[channel_id]
	}
}

/**
 * subscribe to this channel
 */
func (c *Channel) Subscribe(address net.Addr) {
	if subscriber, ok := c.Addresses[address.String()]; ok {
		subscriber.last_tx = time.Now()
	} else {
		c.LastKey += 1
		subscriber := Subscriber{c.LastKey, address, time.Now()}
		c.Subscribers[c.LastKey] = subscriber
		c.Addresses[address.String()] = subscriber
	}
}

func (c *Channel) Unsubscribe(address net.Addr) {
	if subscriber, ok := c.Addresses[address.String()]; ok {
		delete(c.Subscribers, subscriber.key)
		delete(c.Addresses, subscriber.address.String())
	}
}

// broadcast this carrier event to all on this channel
// and always return to sender (who can ignore if they wish, or can use as net sidetone
func (c *Channel) Broadcast(event bitoip.CarrierEventPayload, localAddress *net.UDPAddr) {
	for _, v := range c.Subscribers {
		// don't broadcast back to sender
		bitoip.UDPTx(bitoip.CarrierEvent, event, v.address.String(), localAddress)
	}

}


