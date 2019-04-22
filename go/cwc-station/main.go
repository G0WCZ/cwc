package main

import (
	"../cwc"
	"context"
	"flag"
	"fmt"
	"log"
)

const maxBufferSize = 508
const CqPort = 5990
const LocalMulticast = "224.0.0.73:%d"

func main() {
	var refAddPtr = flag.String("ref", "cwc0.nodestone.io:7388", "--ref=host:port")
	var cqPtr = flag.Bool("cq", false, "--cq is CQ mode, no server, local broadcast")
	var localPort = flag.Int("port", CqPort, "--port=<local-udp-port>")
	var keyIn = flag.String("keyin", "17", "-keyin=17")
	var keyOut = flag.String("keyout", "27", "-keyout=27")
	var serialDevice = flag.String("serial", "", "-serial=<serial-device-name>")
	var testFeeback = flag.Bool("test", false, "--test to put into local feedback test")

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
	} else {
		morseIO =  cwc.NewPiGPIO()
	}
	morseIO.SetConfig("keyIn", *keyIn)
	morseIO.SetConfig("keyOut", *keyOut)

	if cqMode {
		mcAddress := fmt.Sprintf(LocalMulticast, *localPort)
		log.Printf("Starting in CQ mode with local multicast address %s", mcAddress)

		cwc.StationClient(ctx, true, mcAddress, morseIO, *testFeeback)
	} else {
		log.Printf("Connecting to reflector %s", refAddress)

		cwc.StationClient(ctx, false, refAddress, morseIO, *testFeeback)
	}
}
