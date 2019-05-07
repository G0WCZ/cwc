package main

import "fmt"
import "../bitoip"

/*
 * Protocol Version using semantic versioning
 * See: https://semver.org/
 */
var stationVersion = "1.0.1"

func ReflectorVersion() string {
	return stationVersion
}

func DisplayVersion() string {
	return fmt.Sprintf("CWC Station %s / Protocol %s", stationVersion, bitoip.ProtocolVersionString())
}