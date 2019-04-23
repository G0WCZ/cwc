package bitoip

import (
	"testing"
		"gotest.tools/assert"
	)

func TestTimeSyncMessageEncode(t *testing.T) {
	verb := TimeSync
	payload := TimeSyncPayload{1}

	myBytes := EncodePayload(verb, payload)
	assert.DeepEqual(t, myBytes, []uint8{146, 0, 0, 0, 0, 0, 0, 0, 1})
}

func TestEnumerateChannels(t *testing.T) {
	verb := EnumerateChannels
	payload := Payload(nil)

	myBytes := EncodePayload(verb, payload)
	assert.DeepEqual(t, myBytes, []uint8{EnumerateChannels})
	newVerb, _ := DecodePacket(myBytes)

	assert.Equal(t, newVerb, EnumerateChannels)}

func TestListChannels(t *testing.T) {
	verb := ListChannels
	payload := ListChannelsPayload{[MaxChannelsPerMessage]uint16{1, 3, 2}}

	myBytes := EncodePayload(verb, payload)
	assert.DeepEqual(t, myBytes[0:7], []uint8{ListChannels, 0, 1, 0, 3, 0, 2})

	newVerb, newInt := DecodePacket(myBytes)

	lcp := newInt.(*ListChannelsPayload)
	assert.DeepEqual(t, lcp.Channels[0:3], []uint16{0x01, 0x03, 0x02})
	assert.Equal(t, newVerb, ListChannels)
}

func TestTimeSyncResponse(t *testing.T) {
	verb := TimeSyncResponse
	payload := TimeSyncResponsePayload{
		1,
		2}

	myBytes := EncodePayload(verb, payload)
	assert.DeepEqual(t, myBytes, []uint8{TimeSyncResponse,
		0, 0, 0, 0, 0, 0, 0, 1,
		0, 0, 0, 0, 0, 0, 0, 2})

	newVerb, newInt := DecodePacket(myBytes)

	tsrp := newInt.(*TimeSyncResponsePayload)
	assert.Equal(t, tsrp.GivenTime, int64(0x01))
	assert.Equal(t, newVerb, TimeSyncResponse)}

func TestListenRequestPayload (t *testing.T) {
	verb := ListenRequest
	var callsign Callsign
	copy(callsign[:], []byte("G0WCZ"))

	payload := ListenRequestPayload{ 99, callsign}

	myBytes := EncodePayload(verb, payload)

	assert.DeepEqual(t, myBytes[0:8], []uint8{ListenRequest, 0, 99,
										0x47, 0x30, 0x57, 0x43, 0x5a})

	newVerb, newInt := DecodePacket(myBytes)

	lcp := newInt.(*ListenRequestPayload)
	assert.Equal(t, lcp.Channel, uint16(99))
	assert.Equal(t, newVerb, ListenRequest)
}

func TestListenConfirmPayload (t *testing.T) {
	verb := ListenConfirm

	payload := ListenConfirmPayload{ 99, 0xeeee}

	myBytes := EncodePayload(verb, payload)

	assert.DeepEqual(t, myBytes, []uint8{ListenConfirm, 0, 99, 0xee, 0xee})

	newVerb, newInt := DecodePacket(myBytes)

	lcp := newInt.(*ListenConfirmPayload)
	assert.Equal(t, lcp.CarrierKey, uint16(0xeeee))
	assert.Equal(t, newVerb, ListenConfirm)}

func TestUnlistenPayload (t *testing.T) {
	verb := Unlisten

	payload := UnlistenPayload{ 99, 0xeeee}

	myBytes := EncodePayload(verb, payload)

	assert.DeepEqual(t, myBytes, []uint8{Unlisten, 0, 99, 0xee, 0xee})

	newVerb, newInt := DecodePacket(myBytes)

	ulp := newInt.(*UnlistenPayload)
	assert.Equal(t, ulp.CarrierKey, uint16(0xeeee))
	assert.Equal(t, newVerb, Unlisten)
}

func TestKeyValuePayload (t *testing.T) {
	verb := KeyValue

	var key [8]byte
	var value [16]byte

	copy(key[:], []byte("keyx"))
	copy(value[:], []byte("somevalue"))

	payload := KeyValuePayload{
		99, 0xeeee, key, value }

	myBytes := EncodePayload(verb, payload)

	assert.DeepEqual(t, myBytes, []uint8{KeyValue,
		0, 99, 0xee, 0xee,
		0x6b, 0x65, 0x79, 0x78, 0x00, 0x00, 0x00, 0x00,
		0x73, 0x6f, 0x6d, 0x65, 0x76, 0x61, 0x6c, 0x75,
		0x65, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	newVerb, newInt := DecodePacket(myBytes)

	kvp := newInt.(*KeyValuePayload)
	assert.Equal(t, kvp.Channel, uint16(99))
	assert.Equal(t, newVerb, KeyValue)}


func TestCarrierEventPayload (t *testing.T) {
	verb := CarrierEvent

	onEvent := CarrierBitEvent{0, BitOn}
	offEvent := CarrierBitEvent{100, BitOff}
	lastEvent := CarrierBitEvent{100, LastEvent}

	payload := CarrierEventPayload{
		99, 0xeeee,
		0,
		[MaxBitEvents]CarrierBitEvent{onEvent, offEvent, lastEvent},
		 0}

	myBytes := EncodePayload(verb, payload)

	assert.DeepEqual(t, myBytes[0:28], []uint8{
		CarrierEvent, 0, 99, 0xee, 0xee,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, uint8(BitOn),
		0, 0, 0, 100, uint8(BitOff),
		0, 0, 0, 100, uint8(LastEvent),
	})

	newVerb, newInt := DecodePacket(myBytes)

	cep := newInt.(*CarrierEventPayload)
	assert.Equal(t, cep.Channel, uint16(99))
	assert.Equal(t, newVerb, CarrierEvent)
}