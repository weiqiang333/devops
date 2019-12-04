package web

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/spf13/viper"

	"github.com/weiqiang333/devops/internal/crontab"
	"github.com/weiqiang333/devops/web/handlers/auth"
	"github.com/weiqiang333/devops/web/handlers/service"
)


func Web()  {
	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()
	// Logging to a file.
	f, _ := os.Create("logs/access.log")
	gin.DefaultWriter = io.MultiWriter(f)

	router := gin.New()
	router.LoadHTMLGlob("web/templates/*")
	router.Static("/static", "web/static")

	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())
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
