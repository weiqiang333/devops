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
	"github.com/weiqiang333/devops/internal/release/pre_release"
)


//PostPreRelease 调用 jenkins 钩子触发编译
func PostPreRelease(c *gin.Context)  {
	username := fmt.Sprint(auth.Me(c))
	jobs := c.PostFormArray("job")
	releaseNote := c.PostForm("RELEASE_NOTE")
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
	status := pre_release.TriggerBuildJobs(username, releaseJobs, releaseNote)
	c.JSON(http.StatusOK, gin.H{
		"response": status,
	})
	return
}


//GetPreRelease Handle
func GetPreRelease(c *gin.Context)  {
	username := auth.Me(c)
	releaseJobs, err := pre_release.GetJobs("")
	releaseJobsBuilds, _ := pre_release.GetBuilds("build")
	if err != nil {
		c.HTML(http.StatusNotImplemented, "release/pre-release.tmpl", gin.H{
			"user": username,
			"release": "action",
			"error": fmt.Sprintf("jobs 获取异常 %v", err),
		})
		return
	}
	c.HTML(http.StatusOK, "release/pre-release.tmpl", gin.H{
		"user": username,
		"release": "action",
		"releaseJobs": releaseJobs,
		"releaseJobsBuilds": releaseJobsBuilds,
	})
	return
}
