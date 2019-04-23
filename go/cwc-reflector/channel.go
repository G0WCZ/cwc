package main

import (
	"log"
	"net"
)
import (
	"../bitoip"
	"time"
)

type Subscriber struct {
	key bitoip.CarrierKeyType
	address net.UDPAddr
	last_tx time.Time
}

type Channel struct {
	ChannelId bitoip.ChannelIdType
	Subscribers map[bitoip.CarrierKeyType]Subscriber
	Addresses map[string]Subscriber
	LastKey bitoip.CarrierKeyType
}

var channels = make(map[uint16]*Channel)

func NewChannel(channelId bitoip.ChannelIdType) Channel {
	return Channel {
		 channelId,
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
func GetChannel(channel_id bitoip.ChannelIdType) *Channel {
	if channel, ok := channels[channel_id]; ok {
		return channel
	} else {
		nc := NewChannel(channel_id)
		channels[channel_id] = &nc
		return &nc
	}
}

/**
 * subscribe to this channel
 */
func (c *Channel) Subscribe(address net.UDPAddr) bitoip.CarrierKeyType {
	log.Printf("subscribe from: %v", address)
	log.Printf("channels: %v", channels)
	if subscriber, ok := c.Addresses[address.String()]; ok {
		subscriber.last_tx = time.Now()
		log.Printf("suscribe existing key %d", subscriber.key)
		return subscriber.key
	} else {
		c.LastKey += 1
		subscriber := Subscriber{c.LastKey, address, time.Now()}
		c.Subscribers[c.LastKey] = subscriber
		c.Addresses[address.String()] = subscriber
		log.Printf("suscribe new key %d", subscriber.key)
		return subscriber.key
	}
}

func (c *Channel) Unsubscribe(address net.UDPAddr) {
	if subscriber, ok := c.Addresses[address.String()]; ok {
		delete(c.Subscribers, subscriber.key)
		delete(c.Addresses, subscriber.address.String())
	}
}

// broadcast this carrier event to all on this channel
// and always return to sender (who can ignore if they wish, or can use as net sidetone
func (c *Channel) Broadcast(event bitoip.CarrierEventPayload) {
	for _, v := range c.Subscribers {
		log.Printf("sending to subs %v: %v", v.address, event)
		bitoip.UDPTx(bitoip.CarrierEvent, event, &v.address)
	}

}


