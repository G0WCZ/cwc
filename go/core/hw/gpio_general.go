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

import (
	"github.com/G0WCZ/cwc/config"
)

type GPIOGeneral struct {
	Config      *config.Config
	adapterName string
}

func (G *GPIOGeneral) Open() error {
	panic("implement me")
}

func (G *GPIOGeneral) ConfigChanged() error {
	panic("implement me")
}

func (G *GPIOGeneral) SetStatus(string, string) {
	panic("implement me")
}

func (G *GPIOGeneral) GetStatus(string) string {
	panic("implement me")
}

func (G *GPIOGeneral) Close() error {
	panic("implement me")
}

func NewGPIOGeneral(config *config.Config, adapterName string) GeneralIO {
	return &GPIOGeneral{Config: config, adapterName: adapterName}
}
