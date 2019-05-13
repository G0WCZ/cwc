package main

import "fmt"
import "../bitoip"

/*
 * Protocol Version using semantic versioning
 * See: https://semver.org/
 */
var reflectorVersion = "2.0.1"

func ReflectorVersion() string {
	return reflectorVersion
}

func DisplayVersion() string {
	return fmt.Sprintf("CWC Reflector %s / Protocol %s", reflectorVersion, bitoip.ProtocolVersionString())
}