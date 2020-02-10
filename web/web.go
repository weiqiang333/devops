package web

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/thinkerou/favicon"

	"github.com/weiqiang333/devops/web/handlers/auth"
	"github.com/weiqiang333/devops/web/handlers/aws/rds"
	"github.com/weiqiang333/devops/web/handlers/ldapadmin"
	"github.com/weiqiang333/devops/web/handlers/service"
	"github.com/weiqiang333/devops/web/handlers/release"
)


func Web()  {
	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()
	// Logging to a file.
	f, _ := os.Create("logs/access.log")
	gin.DefaultWriter = io.MultiWriter(f)

	router := gin.New()
	router.SetFuncMap(template.FuncMap{
		"ifEqual": ifEqual,
		"formatAsDate": formatAsDate,
	})
	router.LoadHTMLGlob("web/templates/**/*")
	router.Static("/static", "web/static")
	router.Use(favicon.New("web/static/images/favicon.png"))

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
	{
		router.GET("/status", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})
		forgetLdap := router.Group("/ldapAdmin")
		{
			forgetLdap.GET("/forgetPwd", auth.GetForgetPwd)
			forgetLdap.POST("/forgetPwd/postMailVerificationCode", auth.PostMailVerificationCode)
			forgetLdap.POST("/forgetPwd/modifyUserPwd", auth.ModifyUserPwd)
		}
	}

	router.Use(sessions.Sessions("mysession", sessions.NewCookieStore([]byte("secret"))))
	{
		router.GET("/", func(c *gin.Context) {
			username := auth.Me(c)
			c.HTML(http.StatusOK, "base/index.tmpl", gin.H{
				"home": "active",
				"user": username,
			})
		})
		router.POST("/login", auth.Login)
		router.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "user/login.tmpl", gin.H{})
		})
		router.GET("/logout", auth.Logout)
	}

	// Private
	private := router.Group("/", auth.AuthRequired)
	{
		private.GET("/users", auth.Users)
		private.GET("/users/createqr", auth.CreateQRcode)

		private.GET("/service", service.ListService)
		private.POST("/service", service.ListService)

		ldap := private.Group("/ldapAdmin")
		{
			ldap.GET("/", ldapadmin.LadpAdmin)
			ldap.GET("/modifyUserPwd", ldapadmin.GetModifyPwdPage)
			ldap.POST("/modifyUserPwd", ldapadmin.ModifyUserPwd)
		}

		aws := private.Group("/awsAdmin")
		{
			aws.GET("/rdsRsyncWorkorder", rds.GetSyncWorkorder)
			aws.POST("/rdsRsyncWorkorder", rds.PostRsyncWorkorder)
			aws.GET("/rdsRsyncWorkorder/:id", rds.GetOrderId)
			aws.POST("/rdsRsyncWorkorder/:id", rds.PostApproval)
			aws.GET("/rdsRsyncWorkorder/:id/rsync", rds.GetRsyncDegree)
			aws.POST("/rdsRsyncWorkorder/:id/rsync", rds.PostExecuteRsync)
		}

		releases := private.Group("/release")
		{
			releases.GET("/", release.GetRelease)
			releases.GET("/pre-release", release.GetPreRelease)
			releases.POST("/pre-release", release.PostPreRelease)
		}
	}

	err := router.Run(viper.GetString("address")) // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Println(err.Error())
	}
}


//ifEqual gin router SetFuncMap
func ifEqual(a, b interface{}) string {
	aa := fmt.Sprint(a)
	bb := fmt.Sprint(b)
	if aa == bb {
		return aa
	} else {
		return ""
	}
}

func formatAsDate(date time.Time, timezone string) string {
	if date.IsZero() {
		return "无"
	}
	if timezone == "utc" {
		date = date.Add(time.Hour * 8)
	}
	return date.Format("2006-01-02 15:04:05")
}
