package web

import (
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/weiqiang333/devops/web/handlers/service"
	"github.com/weiqiang333/devops/internal/crontab"
	"github.com/weiqiang333/devops/web/handlers/auth"
)


func Web()  {
	router := gin.Default()
	router.LoadHTMLGlob("web/templates/*")
	router.Static("/static", "web/static")

	router.Use(sessions.Sessions("mysession", sessions.NewCookieStore([]byte("secret"))))
	{
		router.POST("/login", auth.Login)
		router.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.tmpl", gin.H{})
		})
		router.GET("/logout", auth.Logout)
	}

	router.GET("/", func(c *gin.Context) {
		username := auth.Me(c)
		fmt.Println(username)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"home": "active",
			"user": username,
		})
	})
	router.GET("/service", service.ListService)

	// Private
	private := router.Group("/", auth.AuthRequired)
	{
		private.GET("/status", auth.Status)

		private.POST("/service", service.ListService)
	}


	crontab.CronTab()

	err := router.Run(viper.GetString("address")) // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Println(err.Error())
	}
}
