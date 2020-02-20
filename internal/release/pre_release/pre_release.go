/*
	pre release, 获取 jobs 和 触发 jobs 两种方法。
	用于上线前触发 job 构建操作
 */
package pre_release

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/weiqiang333/devops/internal/database"
	"github.com/weiqiang333/devops/internal/model"
)


//GetJobs 获取 jobs
func GetJobs(job string) ([]model.ReleaseJobs, error) {
	// job 为："'backend', 'frontend', 'accounts'" 或 ""
	sql := fmt.Sprintf(`SELECT * FROM release_jobs;`)
	if job != "" {
		sql = fmt.Sprintf(`SELECT * FROM release_jobs WHERE jobname IN (%s);`, job)
	}
	releaseJobs := []model.ReleaseJobs{}
	db := database.Db()
	defer db.Close()
	row, err := db.Query(sql)
	defer row.Close()
	if err != nil {
		log.Printf("select release_jobs error: %v", err)
		return releaseJobs, fmt.Errorf("select fail")
	}

	for row.Next() {
		releaseJob := model.ReleaseJobs{}
		err = row.Scan(&releaseJob.Id, &releaseJob.JobName, &releaseJob.JobUrl, &releaseJob.JobHook, &releaseJob.UpdatedAt, &releaseJob.LastExecuteAt)
		if err != nil {
			log.Printf("GetJobs Scan fail: %v", err)
			continue
		}
		releaseJobs = append(releaseJobs, releaseJob)
	}
	return releaseJobs, nil
}


//TriggerBuildJobs 调用 Jenkins 钩子, 触发构建
func TriggerBuildJobs(username string, jobs []model.ReleaseJobs, releaseNote string) string {
	status := ""
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	for _, job := range jobs {
		url := fmt.Sprintf("%s&BUILD_USER=%s&RELEASE_NOTE=%s", job.JobHook, username, releaseNote)
		resp, err := client.Get(url)
		status += fmt.Sprintf("%s: %s\n", job.JobName, resp.Status)
		if err != nil {
			log.Printf("release PushJobs fail %s %s: %v", job.JobName, resp.Status, err)
			continue
		}
		updateJobLastExecuteTime(job.JobName)
	}
	return status
}

func updateJobLastExecuteTime(job string)  {
	sql := fmt.Sprintf(`
		UPDATE release_jobs
		SET 
			last_execute_at = now() at time zone 'utc'
		WHERE
			jobname = '%s';`, job)
	db := database.Db()
	defer db.Close()
	row, err := db.Query(sql)
	defer row.Close()
	if err != nil {
		log.Printf("UPDATE release_jobs error for updateJobLastExecuteTime: %s - %v", job, err)
	}
}
