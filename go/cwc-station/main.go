package main

import (
	"../cwc"
	"context"
	"flag"
	"fmt"
	"github.com/golang/glog"
)

const maxBufferSize = 508
const CqPort = 5990
const LocalMulticast = "224.0.0.73:%d"

func main() {
	var refAddPtr = flag.String("ref", "cwc0.nodestone.io:7388", "--ref=host:port")
	var cqPtr = flag.Bool("cq", false, "--cq is CQ mode, no server, local broadcast")
	var localPort = flag.Int("port", CqPort, "--port=<local-udp-port>")
	var keyIn = flag.String("keyin", "17", "-keyin=17 sets BCM gpio pin as morse input")
	var keyOut = flag.String("keyout", "27", "-keyout=27 sets BCM gpio pin as morse out")
	var pcmOut = flag.String("pcmout", "13", "-pcmout=13 sets BCM gpio pin as pwm sound out")
	var serialDevice = flag.String("serial", "", "-serial=<serial-device-name>")
	var testFeedback = flag.Bool("test", false, "--test to put into local feedback test mode")
	var sidetoneFreq = flag.String("sidetone", "0", "-sidetone 450 to send 450hz tone on keyout")
	var echo = flag.Bool("echo", false, "-echo turns on remote echo of all sent morse")
	var channel = flag.Int("ch", 0, "-ch <n> to connect to the channel n")
	var callsign = flag.String("de", "", "-de <callsign>")
	var noIO = flag.Bool("noio",false, "-noio uses fake morse IO connections")

	flag.Parse()

	// Mode and address
	cqMode := *cqPtr
	refAddress := *refAddPtr

	// context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Morse Hardware
	var morseIO cwc.IO

	if len(*serialDevice) > 0 {
		morseIO = cwc.NewSerialIO()
		morseIO.SetConfig("serialDevice", *serialDevice)
	} else if *noIO == true {
		morseIO = cwc.NewNullIO()
	} else {
		morseIO =  cwc.NewPiGPIO()
	}
	morseIO.SetConfig("keyIn", *keyIn)
	morseIO.SetConfig("keyOut", *keyOut)
	morseIO.SetConfig("pcmOut", *pcmOut)
	morseIO.SetConfig("sidetoneFreq", *sidetoneFreq)

	if cqMode {
		mcAddress := fmt.Sprintf(LocalMulticast, *localPort)
		glog.Infof("Starting in CQ mode with local multicast address %s", mcAddress)

		cwc.StationClient(ctx, true, mcAddress, morseIO, *testFeedback, *echo, uint16(*channel), *callsign)
	} else {
		glog.Infof("Connecting to reflector %s", refAddress)

		cwc.StationClient(ctx, false, refAddress, morseIO, *testFeedback, *echo, uint16(*channel), *callsign)
	}
}
