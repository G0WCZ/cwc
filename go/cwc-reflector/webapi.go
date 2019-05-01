package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

const APIPort = ":7380"

func APIServer(ctx context.Context, channels *ChannelMap) {
	router := gin.Default()
	router.GET("/channels", func(c *gin.Context) {
		c.JSON(http.StatusOK, *channels)
	})

	router.Run(APIPort)
}

