# CWC - The CW Commuter

## The Idea

The idea:  A little box that you can plug a key and headphones into.  It Wifi connects to your phone 
hotspot.  It has a channel selector.   Dial up a channel and tx/rx CW on that channel. That's it.    

This is an internet transceiver for CW that you can take with you.  It aims to be more like a radio than a computer. 

## What is in the box
A Raspberry Pi or arduino with WiFi.  A few components for audio out and key in.
There will be a channel knob one day and a signal LED.

## Get Started
It is a bit basic for now, but is simple enough and it works:
1. Raspberry Pi:  Get a PI zero W or a PI 3.  Get it on wifi.  
1. Get latest CWC station binary here: https://github.com/G0WCZ/cwc/releases
1. See the wiki for hardware setup: https://github.com/G0WCZ/cwc/wiki/Pi-Zero---Pi-hardware-setup
log into your Raspberry PI.
1. Run a command where you put that binary:
```
./cwc-station -de <your-call> -sidetone <freq>

# example commands:
# get started with you call as Q1AAA and a 500hz sidetone
./cwc-station -de Q1AAA -sidetone 500

# add verbose log outout
./cwc-station -de Q1AAA -sidetone 500 -logtostderr -v 2

# change PI pins:
./cwc-station -de Q1AAA -sidetone <freq>

# get help on commands
./cwc-station -h
```

This is the simple way to get started with Raspberry PI using GPIO pins.
You can also use a USB serial port on desktops etc.  See the wiki for more details.


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

* early release, still in development: Raspberry Pi GPIO / or Mac & Linux * maybe windows with serial port
* planning for: Arduino/NodeMCU

# Pi Zero default setup

Here are the default: GPIO pinouts:
```
BCM17 - connector pin 11 - morse out
BCM27 - connector pin 13 - morse in - use 10k pullup resistor to 3V3
BCM13 - conncetor pin 33 - PWM audio morse out.  
```
See https://cdn-learn.adafruit.com/downloads/pdf/adding-basic-audio-ouput-to-raspberry-pi-zero.pdf for details
of audio output circuit.  Basically voltage divider and low-pass filter to make a head phone output.

Full set of command-line flags:
```
Usage of ./cwc-station:
  -alsologtostderr
    	log to standard error as well as files
  -ch int
    	-ch <n> to connect to the channel n
  -cq
    	--cq is CQ mode, no server, local broadcast
  -de string
    	-de <callsign>
  -echo
    	-echo turns on remote echo of all sent morse
  -keyin string
    	-keyin=17 sets BCM gpio pin as morse input (default "17")
  -keyout string
    	-keyout=27 sets BCM gpio pin as morse out (default "27")
  -log_backtrace_at value
    	when logging hits line file:N, emit a stack trace
  -log_dir string
    	If non-empty, write log files in this directory
  -logtostderr
    	log to standard error instead of files
  -noio
    	-noio uses fake morse IO connections
  -pcmout string
    	-pcmout=13 sets BCM gpio pin as pwm sound out (default "13")
  -port int
    	--port=<local-udp-port> (default 5990)
  -ref string
    	--ref=host:port (default "cwc0.nodestone.io:7388")
  -serial string
    	-serial=<serial-device-name>
  -sidetone string
    	-sidetone 450 to send 450hz tone on keyout (default "0")
  -stderrthreshold value
    	logs at or above this threshold go to stderr
  -test
    	--test to put into local feedback test mode
  -v value
    	log level for V logs
  -vmodule value
    	comma-separated list of pattern=N settings for file-filtered logging
```

# Who did this
Ideas by Grae G0WCZ and the online radio club MX0ONL

Go implementation (for RPi and others) by Grae G0WCZ




