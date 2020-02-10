package release

import (
	"github.com/gin-gonic/gin"
	"github.com/weiqiang333/devops/web/handlers/auth"
	"net/http"
)


//GetRelease /release url
func GetRelease(c *gin.Context)  {
	username := auth.Me(c)
	c.HTML(http.StatusOK, "release/release.tmpl", gin.H{
		"user": username,
		"release": "action",
	})
}
