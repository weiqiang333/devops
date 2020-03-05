package release

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/weiqiang333/devops/internal/release/release_api"
)


//JenkinsCallbackAPI
func JenkinsCallbackAPI(c *gin.Context)  {
	jobName, nOk := c.GetQuery("jobname")
	jobId, iOk := c.GetQuery("id")
	buildResult, rOk := c.GetQuery("result")
	if ! nOk && ! iOk && ! rOk {
		c.String(http.StatusBadRequest, "参数提供有误")
		return
	}
	go release_api.CallbackGrab(jobName, jobId, buildResult)
	c.String(http.StatusOK, "release api 回调以提交")
	return
}