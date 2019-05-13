package bitoip

import "fmt"

/*
 * Protocol Version using semantic versioning
 * See: https://semver.org/
 */
const MajorVersion = uint8(2)
const MinorVersion = uint8(0)
const PatchVersion = uint8(1)

var protocolVersionBytes = []byte{MajorVersion, MinorVersion, PatchVersion}
var protocolVersionString = fmt.Sprintf("%d.%d.%d", MajorVersion, MinorVersion, PatchVersion)

func ProtocolVersionBytes() []byte {
	return protocolVersionBytes
}

func ProtocolVersionString() string {
	return protocolVersionString
}