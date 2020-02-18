/*
	管理配置, 上线发布, 触发 jenkins 钩子预编译

 */
package release

import (
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/weiqiang333/devops/internal/authentication"
	"github.com/weiqiang333/devops/internal/release/pre_release"
	"github.com/weiqiang333/devops/web/handlers/auth"
)


func GetPreReleaseAdmin(c *gin.Context) {
	username := fmt.Sprint(auth.Me(c))
	if ! authentication.SecurityGroupAuthorization(username) {
		c.String(http.StatusUnauthorized, "请确认您是否拥有权限")
		return
	}
	releaseJobs, err := pre_release.GetJobs("")
	if err != nil {
		c.HTML(http.StatusNotImplemented, "release/pre-release-admin.tmpl", gin.H{
			"user": username,
			"release": "action",
			"error": fmt.Sprintf("jobs 获取异常 %v", err),
		})
		return
	}
	c.HTML(http.StatusOK, "release/pre-release-admin.tmpl", gin.H{
		"user": username,
		"release": "action",
		"releaseJobs": releaseJobs,
	})
	return
}

func PostPreReleaseAdmin(c *gin.Context) {
	username := fmt.Sprint(auth.Me(c))
	if ! authentication.SecurityGroupAuthorization(username) {
		c.String(http.StatusUnauthorized, "请确认您是否拥有权限")
		return
	}
	jobName := c.PostForm("jobName")
	jobUrl := c.PostForm("jobUrl")
	jobHook := c.PostForm("jobHook")
	if len(jobName) == 0 || len(jobUrl) == 0 || len(jobHook) == 0 {
		c.String(http.StatusBadRequest, "请正确提交")
		return
	}

	err := pre_release.InsertReleaseJob(jobName, jobUrl, jobHook)
	if err != nil {
		c.String(http.StatusInternalServerError, "发生内部错误：" + err.Error())
		return
	}
	c.String(http.StatusOK, "提交成功")
	return
}
