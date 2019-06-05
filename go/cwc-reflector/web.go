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
package main

import (
	"context"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"net/http"
)

const APIPort = ":7380"

func renderer() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", "web/tmpl/base.html", "web/tmpl/index.html")
	return r
}

func APIServer(ctx context.Context, channels *ChannelMap, address string) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Static("/static", "./web/root")
	router.HTMLRender = renderer()

	router.GET("/api/channels", func(c *gin.Context) {
		c.JSON(http.StatusOK, *channels)
	})
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index", gin.H{
			"HostAndPort": address,
			"Channels":    channels,
		})
	})

	router.Run(APIPort)
}
