package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"../cwc"
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

		cwc.StationClient(ctx, true, mcAddress, morseIO)
	} else {
		log.Printf("Connecting to reflector %s", refAddress)

		cwc.StationClient(ctx, false, refAddress, morseIO)
	}
}

func ReflectorServer(ctx context.Context, address string) {
	pc, err := net.ListenPacket("udp", address)

	if err != nil {
		return
	}

	defer pc.Close()

	buffer := make([]byte, maxBufferSize)
	doneChan := make(chan error, 1)

	go func() {
		for {
			_, _, err := pc.ReadFrom(buffer)

			if err != nil {
				doneChan <- err
				return
			}
			log.Printf("packet rx: %#v", buffer)
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("cancelled")
		err = ctx.Err()
	case err = <-doneChan:
	}

}