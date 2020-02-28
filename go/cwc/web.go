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
package cwc

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)


const WebPort = ":443"

var cfg *Config
var ver string

func WebServer(ctx context.Context, config *Config, version string) {

	cfg = config
	ver = version
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.LoadHTMLFiles("index.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	router.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})
	router.RunTLS(WebPort, "./certs/server.pem", "./certs/server.key")
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		m := string(msg)

		if m == "fromC:status:connected" {
			var res string
			res = fmt.Sprintf("toC:wpm:%d", cfg.KeyerSpeed)
			conn.WriteMessage(t, []byte(res))
//			res = fmt.Sprintf("toC:ver:%s", StationVersion)
//			conn.WriteMessage(t, []byte(res))
                        res = fmt.Sprintf("toC:chan:%d", cfg.Channel)
                        conn.WriteMessage(t, []byte(res))
		}

		s := strings.Split(m, ":")
		direction, name, value := s[0], s[1], s[2]

		if direction == "toC" {
			return
		}

		if name == "wpm" {
			var wpm int
			wpm, err = strconv.Atoi(value)
			if err == nil {
				cfg.KeyerSpeed = wpm
                                CalcDelays()
			}
		}
                if name == "chan" {
                    var chan64 uint64
                    chan64, err := strconv.ParseUint(value, 16, 16)
                    if err == nil {
                        cfg.Channel = uint16(chan64)
                    }
                }

	}
}

