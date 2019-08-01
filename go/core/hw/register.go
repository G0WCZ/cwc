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

package hw

type InputInfo struct {
	Name string
}

type OutputInfo struct {
	Name string
}

type GeneralInfo struct {
	Name string
}

var Inputs = make(map[string]InputInfo)
var Outputs = make(map[string]OutputInfo)
var GeneralIOs = make(map[string]GeneralInfo)
