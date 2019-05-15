/*
Copyright (C) 2019 Graeme Sutherland, Nodestone Limited


This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package bitoip

import "fmt"

/*
 * Protocol Version using semantic versioning
 * See: https://semver.org/
 */
const MajorVersion = uint8(1)
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