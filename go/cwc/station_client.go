package cwc

// General client
// Can be in CQ mode, in which case all is local muticast on the local network
// Else the client of a reflector
// CQ mode is really simple. Only really have to tx and rx carrier events

func StationClient(cqMode bool, addr string, morseIO IO) {
	// CQ mode
	// listen on mc address
	// look for bit events and send them
	// send using any channel
	// rx all channels
	// that's it

	go RunRx(morseIO)







	// Reflector mode
	// opt: time sync with server
	// opt: set callsign
	// list channels
	// suscribe channel(s)
	// save carrier id


}