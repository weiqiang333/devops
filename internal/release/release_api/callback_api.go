package release_api

import (
	"fmt"
	"log"

	"github.com/bndr/gojenkins"
	"github.com/spf13/viper"

	"github.com/weiqiang333/devops/internal/jenkins"
	"github.com/weiqiang333/devops/internal/database"
)


//CallbackGrab 抓取 job build 信息, 进行记录
//ACTION=build 则记录 pre-release build 信息
func CallbackGrab(jobName, jobId, buildResult string)  {
	log.Printf("release api start Callback Grab: %s/%s", jobName, jobId)
	jenkinsBase := viper.GetString("jenkins.baseurl")
	jenkinsUser := viper.GetString("jenkins.user")
	userToken := viper.GetString("jenkins.token")
	jenkinsClient, err := jenkins.ConnJenkins(jenkinsBase, jenkinsUser, userToken)
	if err != nil {
		log.Printf("CallbackGrab fial: %v", err)
		return
	}
	buildAction, buildEnv, err := getJob(jenkinsClient, jobName, jobId)
	if err != nil {
		log.Printf("CallbackGrab fial: %v", err)
		return
	}
	insertBuild(jobName, jobId, buildResult, buildAction, buildEnv)
	log.Printf("release api done Callback Grab: %s/%s", jobName, jobId)
}

func getJob(jenkinsClient *gojenkins.Jenkins, jobName, jobId string) (string, string, error) {
	job, err := jenkinsClient.GetJob(fmt.Sprintf("%s/%s", jobName, jobId))
	if err != nil {
		return "", "", fmt.Errorf("getJob fail: %v", err)
	}
	parameters := map[string]string{}
	for _, parameter := range job.Raw.Actions[0].Parameters {
		if parameter.Name == "BUILD_USER" {
			parameters["BUILD_USER"] = parameter.Value
		}
		if parameter.Name == "ACTION" {
			parameters["ACTION"] = parameter.Value
		}
		if parameter.Name == "BUILD_ENV" {
			parameters["BUILD_ENV"] = parameter.Value
		}
	}
	if parameters["BUILD_USER"] == "" || parameters["ACTION"] == "" || parameters["BUILD_ENV"] == ""{
		return "", "", fmt.Errorf("parameter's BUILD_USER or ACTION or BUILD_ENV is nil")
	}
	return parameters["ACTION"], parameters["BUILD_ENV"], nil
}

func insertBuild(jobName, jobId, buildResult, buildAction, buildEnv string)  {
	sql := fmt.Sprintf(`
		INSERT INTO release_jobs_builds
			(
				jobname,
				job_id,
				build_result,
				build_action,
				build_env,
				update_at
			)
		VALUES
			(
				'%s',
				%v,
				'%s',
				'%s',
				'%s',
				now() at time zone 'utc'
			)
		ON CONFLICT(jobname, job_id)
		DO UPDATE SET
			build_result=EXCLUDED.build_result,
			build_action=EXCLUDED.build_action,
			build_env=EXCLUDED.build_env,
			update_at=EXCLUDED.update_at`, jobName, jobId, buildResult, buildAction, buildEnv,
	)
	db := database.Db()
	defer db.Close()
	row, err := db.Query(sql)
	defer row.Close()
	if err != nil {
		log.Printf("UPDATE release_jobs_builds error for insertBuild: %s, %s, %s, %s, %s", jobName, jobId, buildResult, buildAction, buildEnv)
	}
}
