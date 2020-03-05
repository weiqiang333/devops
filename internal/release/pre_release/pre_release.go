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
	sql := fmt.Sprintf(`SELECT * FROM release_jobs ORDER BY last_execute_at DESC, jobview, jobname;`)
	if job != "" {
		sql = fmt.Sprintf(`SELECT * FROM release_jobs WHERE jobname IN (%s) ORDER BY last_execute_at DESC, jobview, jobname;`, job)
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
		err = row.Scan(&releaseJob.Id, &releaseJob.JobName, &releaseJob.JobUrl, &releaseJob.JobHook, &releaseJob.UpdatedAt, &releaseJob.LastExecuteAt, &releaseJob.JobView)
		if err != nil {
			log.Printf("GetJobs Scan fail: %v", err)
			continue
		}
		releaseJobs = append(releaseJobs, releaseJob)
	}
	return releaseJobs, nil
}

func GetBuilds(action string) ([]model.ReleaseJobsBuilds, error) {
	ofTime := time.Now().UTC().AddDate(0, 0,-1).Format("2006-01-02 15:04:05-07")
	sql := fmt.Sprintf(`
		SELECT *
		FROM release_jobs_builds
		WHERE update_at = (
			SELECT MAX(update_at)
			FROM release_jobs_builds builds
			WHERE	builds.jobname = release_jobs_builds.jobname)
		AND build_action = '%s'
		AND update_at >= '%s';`, action, ofTime,
	)
	fmt.Println(sql)
	releaseJobsBuilds := []model.ReleaseJobsBuilds{}
	db := database.Db()
	defer db.Close()
	row, err := db.Query(sql)
	defer row.Close()
	if err != nil {
		log.Printf("select release_jobs_builds error: %v", err)
		return releaseJobsBuilds, fmt.Errorf("select fail")
	}

	for row.Next() {
		releaseJobBuild := model.ReleaseJobsBuilds{}
		err = row.Scan(&releaseJobBuild.Id, &releaseJobBuild.JobName, &releaseJobBuild.JobId, &releaseJobBuild.BuildResult,
			&releaseJobBuild.BuildAction, &releaseJobBuild.BuildEnv, &releaseJobBuild.UpdateAt)
		if err != nil {
			log.Printf("GetBuilds Scan fail: %v", err)
			continue
		}
		releaseJobsBuilds = append(releaseJobsBuilds, releaseJobBuild)
	}
	return releaseJobsBuilds, nil
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
