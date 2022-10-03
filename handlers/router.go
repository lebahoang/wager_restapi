package handlers

import (
	"hoang/m/logPkg"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up routes in service
func SetupRoutes() *gin.Engine {
	log = logPkg.GetLog()
	router := gin.New()

	router.POST("/wagers", placeWager)
	router.POST("/buy/:wager_id", buyWager)
	router.GET("/wagers", listWagers)
	router.GET("/ping", ping)
	log.Infoln("routes created")
	return router
}
