package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weiqiang333/devops/web/handlers/service"
)


func Web()  {
	router := gin.Default()
	router.LoadHTMLGlob("web/templates/*")
	router.Static("/static", "web/static")

	router.GET("/", 	func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
		})
	})
	router.GET("/service", service.ListService)
	router.POST("/service", service.ListService)
	router.Run()
}