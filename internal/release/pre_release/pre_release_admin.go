package pre_release

import (
	"log"
	"fmt"

	"github.com/weiqiang333/devops/internal/database"
)


//InsertReleaseJob
func InsertReleaseJob(jobName, jobUrl, jobHook, jobView string) error {
	sql := fmt.Sprintf(`
		INSERT INTO release_jobs
			(jobname, joburl, jobhook, updated_at, jobview)
		VALUES ('%s', '%s', '%s', now() at time zone 'utc', '%s')
		ON CONFLICT (jobname)
		DO UPDATE SET
			joburl = EXCLUDED.joburl,
			jobhook = EXCLUDED.jobhook,
			updated_at = EXCLUDED.updated_at,
			jobview = EXCLUDED.jobview;
	`, jobName, jobUrl, jobHook, jobView)
	db := database.Db()
	defer db.Close()
	row, err := db.Query(sql)
	defer row.Close()
	if err != nil {
		log.Printf("INSERT release_jobs error for InsertReleaseJob: %s - %v", jobName, err)
		return fmt.Errorf("INSERT release_jobs error for InsertReleaseJob: %s - %v", jobName, err)
	}
	return nil
}
