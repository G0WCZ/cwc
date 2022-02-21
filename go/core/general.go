/*
Copyright 2022 Graeme Sutherland


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
package core

import (
	"context"
	"github.com/G0WCZ/cwc/config"
	"github.com/G0WCZ/cwc/core/hw"
	"github.com/golang/glog"
)

var generals []hw.GeneralIO

func General(ctx context.Context, config *config.Config) {

	OpenGenerals(config)

	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func OpenGenerals(config *config.Config) {
	generals = hw.ParseGeneralIOs(config)
	for _, g := range generals {
		glog.Infof("opening general %s", g.Name())
		g.Open()
	}
}

func CloseGenerals(config *config.Config) {
	for _, g := range generals {
		glog.Infof("closing general %s", g.Name())
		g.Close()
	}
}

func SetStatus(name string, value string) {
	for _, g := range generals {
		g.SetStatus(name, value)
	}
}

func GetStatus(name string) []string {
	var vals []string
	for _, g := range generals {
		vals = append(vals, g.GetStatus(name))
	}
	return vals
}
