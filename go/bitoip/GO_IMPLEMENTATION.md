# bitoip: golang implementation

## Basics
This is the go implementation of the protocol messages and a UDP receiver and transmitter
which include marshalling and unmarshalling of messages.

## To be fixed later
The transmission of carrier messages is at their full size. These ought to be 
variable length on the wire.  This should be fixed in a backward-compatible kind of way.

## License
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
