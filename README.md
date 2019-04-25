# CWC - The CW Commuter

## The Idea

The idea:  A little box that you can plug a key and headphones into.  It Wifi connects to your phone 
hotspot.  It has a channel selector.   Dial up a channel and tx/rx CW on that channel. That's it.    

This is an internet transceiver for CW that you can take with you.  It aims to be more like a radio than a computer. 

## What is in the box
A Raspberry Pi or arduino with WiFi.  A few components for audio out and key in.
There will be a channel knob one day and a signal LED.

## Communications
There's a protocol based on UDP packets that sends on and off events.
So if you use the key, you are sending on and off events in UDP packets.

At the receiving end there's something that turns packetised on-offs back into contact closures or a tone in your ears.  

UDP is lossy, so it is more radio-like in that sense.    You might lose some packets,  some QSB shrug.

# Broadcast or Reflector
There are two basic modes.  Your CWC station can broadcast on the local network, or talk to a reflector.

In broadcast mode, UDP multicast is used on the local network.  This is a simplified mode for co-located CW training
or similar.

In reflector mode, the station connects to a central reflector that reflects traffic to other connected stations.

See bitoip.md for the on-the-wire protocol details.

# Implementations

* in development: Raspberry Pi GPIO / or Mac & Linux * maybe windows with serial port
* planning for: Arduino/NodeMCU

# Who did this
Ideas by Grae G0WCZ and the online radio club MX0ONL

Go implementation (for RPi and others) by Grae G0WCZ


