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

/**
 * General types and constants for versions
 */

type RelType uint8

const POC RelType = RelType(0)
const Alpha RelType = RelType(1)
const Beta RelType = RelType(2)
const RC RelType = RelType(3)
const Final RelType = RelType(4)

var RelTypeMap = map[RelType]string{
	POC:   "-pre-alpha",
	Alpha: "-alpha",
	Beta:  "-beta",
	RC:    "-rc",
	Final: "",
}

type Version struct {
	Major   uint8
	Minor   uint8
	Patch   uint8
	Release RelType
}

func (v *Version) Bytes() []byte {
	return []byte{v.Major, v.Minor, v.Patch, uint8(v.Release)}
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d%s", v.Major, v.Minor, v.Patch, RelTypeMap[v.Release])
}
