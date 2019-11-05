package web

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/weiqiang333/devops/web/handlers/service"
	"github.com/weiqiang333/devops/internal/crontab"
)


func Web()  {
	router := gin.Default()
	router.LoadHTMLGlob("web/templates/*")
	router.Static("/static", "web/static")

	router.GET("/", 	func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"home": "active",
		})
	})
	router.GET("/service", service.ListService)
	router.POST("/service", service.ListService)
	crontab.CronTab()
	err := router.Run(viper.GetString("address")) // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Println(err.Error())
	}
}