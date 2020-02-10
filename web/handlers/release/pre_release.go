/*
	上线发布，触发 jenkins 钩子预编译
*/
package release

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/weiqiang333/devops/web/handlers/auth"
	pre_release "github.com/weiqiang333/devops/internal/release/pre-release"
)


//PostPreRelease 将掉 jenkins 钩子触发编译
func PostPreRelease(c *gin.Context)  {
	username := fmt.Sprint(auth.Me(c))
	jobs := c.PostFormArray("job")
	log.Printf("%s PostPreRelease 提交将执行预编译操作：%v", username, jobs)
	if len(jobs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"response": "请选择提交",
		})
		return
	}
	for i, v := range jobs {
		jobs[i] = fmt.Sprintf("'%s'", v)
	}
	job := strings.Join(jobs, ",")
	releaseJobs, err := pre_release.GetJobs(job)
	if err != nil {
		c.JSON(http.StatusNotImplemented, gin.H{
			"response": err.Error(),
		})
		return
	}
	status := pre_release.PushJobs(username, releaseJobs)
	c.JSON(http.StatusOK, gin.H{
		"response": status,
	})
}


//GetPreRelease Handle
func GetPreRelease(c *gin.Context)  {
	username := auth.Me(c)
	releaseJobs, err := pre_release.GetJobs("")
	if err != nil {
		c.HTML(http.StatusNotImplemented, "release/pre-release.tmpl", gin.H{
			"user": username,
			"release": "action",
			"error": fmt.Sprintf("jobs 获取异常 %v", err),
		})
	}
	c.HTML(http.StatusOK, "release/pre-release.tmpl", gin.H{
		"user": username,
		"release": "action",
		"releaseJobs": releaseJobs,
	})
}