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
		return releaseJobs, fmt.Errorf("seleact fail")
	}

	for row.Next() {
		releaseJob := model.ReleaseJobs{}
		err = row.Scan(&releaseJob.Id, &releaseJob.JobName, &releaseJob.JobUrl, &releaseJob.JobHook, &releaseJob.UpdatedAt, &releaseJob.LastExecuteAt)
		if err != nil {
			log.Printf("GetJobs Scan fail: %v", err)
		}
		releaseJobs = append(releaseJobs, releaseJob)
	}
	return releaseJobs, nil
}


//PushJobs
func PushJobs(username string, jobs []model.ReleaseJobs) string {
	status := ""
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	for _, job := range jobs {
		url := fmt.Sprintf("%s&BUILD_USER=%s", job.JobHook, username)
		resp, err := client.Get(url)
		status += fmt.Sprintf("%s: %s\n", job.JobName, resp.Status)
		if err != nil {
			log.Printf("release PushJobs fail %s %s: %v", url, resp.Status, err)
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
		log.Printf("UPDATE release_jobs error for updateJobLastExecuteTime: %s - %v", err)
	}
}
