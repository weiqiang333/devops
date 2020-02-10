package crontab

import (
	"log"

	"github.com/spf13/viper"
	"github.com/robfig/cron"
)


// CronTab is cron task
func CronTab()  {
	updateServerListSpec := viper.GetString("cron.update_server_list")
	updateServiceSpec := viper.GetString("cron.update_service")
	updatePwdExpiredSpec := viper.GetString("cron.update_pwd_expired")
	alarmPwdExpiredSpec := viper.GetString("cron.alarm_pwd_expired")

	c := cron.New()
	err := c.AddFunc(updateServerListSpec, UpdateServerList)
	if err != nil {
		log.Printf("cron run UpdateServerList error %s:", err.Error())
	}

	err = c.AddFunc(updateServiceSpec, UpdateService)
	if err != nil {
		log.Printf("cron run UpdateService error %s:", err.Error())
	}

	err = c.AddFunc(updatePwdExpiredSpec, updatePwdExpired)
	if err != nil {
		log.Printf("cron run updatePwdExpired error %s:", err.Error())
	}

	err = c.AddFunc(alarmPwdExpiredSpec, alarmPwdExpired)
	if err != nil {
		log.Printf("cron run alarmPwdExpired error %s:", err.Error())
	}

	c.Start()
	select{}
}
