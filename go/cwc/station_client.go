package cwc

import (
	"../bitoip"
	"context"
	"github.com/golang/glog"
	"net"
	"strings"
	"time"
)

// General client
// Can be in CQ mode, in which case all is local muticast on the local network
// Else the client of a reflector
// CQ mode is really simple. Only really have to tx and rx carrier events

func StationClient(ctx context.Context, cqMode bool,
	addr string, morseIO IO, testFeedback bool, echo bool,
	channel bitoip.ChannelIdType, callsign string) {

	resolvedAddress, err := net.ResolveUDPAddr("udp", addr)

	if err != nil {
		glog.Errorf("Error resolving address %s %v", addr, err)
		return
	}
	toSend := make(chan bitoip.CarrierEventPayload)
	toMorse := make(chan bitoip.RxMSG)

	// Morse receiver
	go RunMorseRx(ctx, morseIO, toSend, echo, channel)

	localRxAddress, err := net.ResolveUDPAddr("udp", "0.0.0.0:0")

	if err != nil {
		glog.Fatalf("Can't allocate local address: %v", err)
	}

	// UDP Receiver
	go bitoip.UDPRx(ctx, localRxAddress, toMorse)

	if !cqMode {
		time.Sleep(time.Second * 1)
		// TODO: full reflector mode implementation
		// Reflector mode setup
		// 1/ time sync with server
		// 2/ set callsign
		// 3/ list channels
		// 4/ suscribe channel(s)
		// 5/ save carrier id

		var csBase [16]byte
		r := strings.NewReader(callsign)
		_, err := r.Read(csBase[0:16])

		if err != nil {
			glog.Errorf("Callsign %s can not be encoded", callsign)
		}

		bitoip.UDPTx(bitoip.ListenRequest, bitoip.ListenRequestPayload{
			channel,
			csBase,
		},
			resolvedAddress,
		)
	}

	// Do time sync
	const timeOffsetBucketSize = 5

	timeOffsetIndex := 0
	timeOffsets := make([]int64, timeOffsetBucketSize, timeOffsetBucketSize)
	timeOffsetSum := int64(0)

	commonTimeOffset := int64(0)

	bitoip.UDPTx(bitoip.TimeSync, bitoip.TimeSyncPayload{
		time.Now().UnixNano(),
	}, resolvedAddress)

	lastUDPSend := time.Now()

	keepAliveTick := time.Tick(20 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return

		case cep := <-toSend:
			glog.V(2).Infof("carrier event payload to send: %v", cep)
			// TODO fill in some channel details
			bitoip.UDPTx(bitoip.CarrierEvent, cep, resolvedAddress)
			if testFeedback {
				QueueForTransmit(&cep)
			}

		case tm := <-toMorse:
			switch tm.Verb {
			case bitoip.CarrierEvent:
				glog.V(2).Infof("carrier events to morse: %v", tm)
				QueueForTransmit(tm.Payload.(*bitoip.CarrierEventPayload))

			case bitoip.ListenConfirm:
				glog.V(2).Infof("listen confirm: %v", tm)
				lc := tm.Payload.(*bitoip.ListenConfirmPayload)
				glog.Infof("listening channel %d with carrier key %d", lc.Channel, lc.CarrierKey)
				SetCarrierKey(lc.CarrierKey)

			case bitoip.TimeSyncResponse:
				glog.V(2).Infof("time sync response %v", tm)
				tsr := tm.Payload.(*bitoip.TimeSyncResponsePayload)
				now := time.Now().UnixNano()
				latestTimeOffset := ((2*tsr.CurrentTime) - tsr.GivenTime - now) / 2
				timeOffsetSum -= timeOffsets[timeOffsetIndex]
				timeOffsets[timeOffsetIndex] = latestTimeOffset
				timeOffsetSum += latestTimeOffset
				timeOffsetIndex = (timeOffsetIndex + 1) % timeOffsetBucketSize
				commonTimeOffset = timeOffsetSum / timeOffsetBucketSize
				SetTimeOffset(commonTimeOffset)

				glog.V(2).Infof("time offset %d", commonTimeOffset)
			}

		case kat := <-keepAliveTick:
			if kat.Sub(lastUDPSend) > time.Duration(20*time.Second) {
				lastUDPSend = kat
				p := bitoip.CarrierEventPayload{
					channel,
					CarrierKey(),
					time.Now().UnixNano(),
					[bitoip.MaxBitEvents]bitoip.CarrierBitEvent{
						bitoip.CarrierBitEvent{0, bitoip.BitOff | bitoip.LastEvent},
					},
					kat.UnixNano(),
				}
				glog.V(2).Info("sending keepalive")
				bitoip.UDPTx(bitoip.CarrierEvent, p, resolvedAddress)
			}
		}

	}
}
