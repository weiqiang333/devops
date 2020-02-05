package crontab

import (
	"log"

	"github.com/spf13/viper"
	"github.com/robfig/cron"
)


// CronTab is cron task
func CronTab()  {
	//updateCmdbServerSpec := viper.GetString("cron.update_cmdb_server")
	//updateServiceSpec := viper.GetString("cron.update_service")
	//updatePwdExpiredSpec := viper.GetString("cron.update_pwd_expired")
	updateAwsInstanceTypesSpec := viper.GetString("cron.update_aws_inastance_types")
	c := cron.New()
	err := c.AddFunc(updateAwsInstanceTypesSpec, UpdateAwsInstanceTypes)
	if err != nil {
		log.Printf("cron run UpdateAwsInstanceTypes error %s:", err.Error())
	}

	//err = c.AddFunc(updateCmdbServerSpec, UpdateCmdbServer)
	//if err != nil {
	//	log.Printf("cron run UpdateServerList error %s:", err.Error())
	//}
	//
	//err = c.AddFunc(updateServiceSpec, UpdateService)
	//if err != nil {
	//	log.Printf("cron run UpdateService error %s:", err.Error())
	//}
	//err = c.AddFunc(updatePwdExpiredSpec, updatePwdExpired)
	//if err != nil {
	//	log.Printf("cron run updatePwdExpired error %s:", err.Error())
	//}
	c.Start()
}
