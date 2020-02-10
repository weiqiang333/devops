package release

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"github.com/weiqiang333/devops/web/handlers/auth"
)


//GetRelease /release url
func GetRelease(c *gin.Context)  {
	username := auth.Me(c)
	c.HTML(http.StatusOK, "release/release.tmpl", gin.H{
		"user": username,
		"release": "action",
	})
}
