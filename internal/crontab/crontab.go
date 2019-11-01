package crontab

import (
	"github.com/spf13/viper"
	"log"

	"github.com/robfig/cron"
)


// CronTab is cron task
func CronTab()  {
	updateServerListSpec := viper.GetString("cron.update_server_list")
	updateServiceSpec := viper.GetString("cron.update_service")
	c := cron.New()
	err := c.AddFunc(updateServerListSpec, UpdateServerList)
	if err != nil {
		log.Printf("cron run UpdateServerList error %s:", err.Error())
	}

	err = c.AddFunc(updateServiceSpec, UpdateService)
	if err != nil {
		log.Printf("cron run UpdateService error %s:", err.Error())
	}
	c.Start()
}


