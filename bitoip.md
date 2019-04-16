# bitoip
bitoip is the on-the-wire protocol for sending on/off events around.   We use it for the CW Commuter.

## What

bit-over-ip: Low level byte-oriented protocol for sending timed on-off state over IP.

## Why

I'm looking for a way to transmit hand-keyed and keyer-keyed morse (telegraph-style) over the
internet or maybe even over other amateur IP services.

Ideally, this is via a toggled line on a serial port or software generated from morse keying or morse (tone) detecting software.

## The model in brief

A *reflector* can host one or more *channels*, each of which can contain multiple *carriers*.  A *station* can exchange
messages with a *reflector* to interact with its channels. 

This becomes generated *carrier* that is a bunch of on and off events, 
with strict timing maintained so it can be re-constructed at
the far end.  Packets of on/off events are sent via UDP, as is this whole protocol.

Data can get lost, but we can deal with that.

Carriers are grouped into a *channel* which has a name, and some control interface.

Channels are hosted on a *reflector*.

##  Entities in detail

Here is some more detail on what the logical parts are:

A **carrier** is a set of bit on/off events for a single bit.  A carrier has:
 - an id (which is used to differentiate different carrier in a channel)
 - a series of on and off events and control messages
 - some sort of sender id (like a callsign)

Note that a carrier should always end up off, not stuck on.  So timed on/off events must always end with an off.  That means the way we collect and send bit data has to
take that into account.

A **channel** is a named grouping of carriers.  A server may support one or more channels.  A channel has:
 - a control interface to allow:
    - enumeration of carriers
    - subscription and unsubscription
 - a name.  Should be unique in a server
 - combining with a server IP address and port, should be possible to make a channel URI

 Given that a channel can have multiple carriers, the receiving end will need to differentiate them when making them audible by
 using tone offsets or something similar, so indvidual signals can be recognised. Note that a channel originates at a reflector.

A **reflector** can receive and publish one or more channels by communicating with stations and (later) possibly other reflectors.
There is no reflector-to-station routing mechanism proposed here yet.

A **station** can register interest in channels with a reflector and add carriers to a channel.

## How this works

A station would connect to a reflector and get an enumeration of channels.  The station can then
subscribe to a channel to receive carrier packets related to that channel.  The subscription results in time
sync being established for that channel, and then the station will receive packets relating to that channel.  The station
can also transmit packets relating to that channel based on the time sync.

There can be multiple carriers per channel, so it is up to the station to make sense of these.

A station ought to also have a promiscuous mode that will listen to all channels it can hear.  

So the role of a reflector is to:
* answer and process incoming requests for Channel Lists, Time Sync, Listen and Unlisten requests.
* Receive and send KeyValue pairs as needed
* Receive and broadcast incoming Carrier Events to listening stations

## Protocol

This protocol is designed to work okay with UDP or similar lossy protocol.

### Channel access

The protocol is composed of packets containing timed on and off information for a bit, with the
necessary time offsets from the start of the stream.

There is a facility to communicate basic key/value ascii pairs.

Generally, top bit set means that is a control value.

multi-byte values are big-endian

## Enumerate channels

Enumerate channels at this reflector

EN (Enumerate) == 0x90

## Enumerate list

Response to enumerate channels available, responding with a list of channel numbers.

EL (Enumerate List) == 0x91

0x91, channel_no (2 bytes), channel_no (2 bytes), ... , 0x00, 0x00

## Server time sync

There are a pair of time sync events to allow stations and reflectors to calculate lateny in the circuit and adjust for
it if needed.

### Initiate time sync
0x92 current_time [timestamp size]

### Respond to time sync
0x93 given_time_or 0 (from prev 0x91) [8 byte, unix nanoseconds], current_time [8byte unix nanoseconds]


## Listen to channel
LI (Listen) == 0x94

0x94, channel_no (2 bytes), callsign/id (utf8 string), 0x00, [ pass_token, 0x00 ]

## Carrier key return or listen confirm
Send from a reflector to provide carrier key after listening

0x95, channel_no (2 byte), carrier_key (2 bytes)

## Unlisten
UL (unlisten) = 0x96

0x96, channel_no (2 bytes), carrier_key (2 bytes)


### Key-value pairs
KV (Key Value) == 0x81:

0x80, channel_no (2 bytes), carrier_key (2 bytes), key(utf-8 string), 0x00, value(utf-8 string), 0x00

Used key-value pairs:
DE callsign of sender
IN general info


### carrier events

BE (Bit Event: timed) == 0x82:

0x82, channel_key (2 bytes), carrier_key (2 bytes), start-8-byte-timestamp (ns), [(time-offset-ns (4 byte), event_type), ...]

event_type is:
0x00: bit off
0x01: bit on

flags: bitwise flags
Not currently used

## Stream Semantics -> sample interaction
```
# basic setup of chnnel
(optional) station sends 0x92 timesync
(optional) reflector sends  0x93 timesync response
station sends 0x90 enumerate channels
reflector sends 0x91 channel info responses

# listen to channel
station sends 0x94 listen to channel
reflector sends back 0x95 carrier key

# send/recieve channel info
station sends 0x82 bit events
reflector sends 0x82 bit events
...

# unlisten
station sends 0x96 unlisten to channel with channel key

```


## Example converstation 

### Initial time sync
```
station: 92 15 44 64 52 e0 a9 50 05
reflector: 93 15 44 64 52 e0 a9 50 05 15 44 64 52 e0 a9 90 05
```

## Enumerate channels
```
station: 90
reflector: 91 00 01 00 02 00 11 00 00
reflector: 91 ea bb cc dd ee ee 00 00
```

## listen to a channel
```
station: 94 00 01 'G' '0' 'W' 'C' 'Z' 00
reflector: 95 00 01 00 02
```
# send a dit at 25 wpm
```
station: 15 00 01 00 02 44 64 52 e0 c9 50 05 00 00 00 00 01 00 00 05 00 00
```

# unlisten channel
```
reflector: 96 00 01 00 02
```

