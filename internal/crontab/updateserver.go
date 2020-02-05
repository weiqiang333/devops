package crontab

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/weiqiang333/devops/internal/awscli"
	"github.com/weiqiang333/devops/internal/database"
	"github.com/weiqiang333/devops/internal/model"
)


func readsAWSServer() []model.Instance {
	var instances = []model.Instance{}
	svc := awscli.AwsCli()
	input := &ec2.DescribeInstancesInput{}
	result, err := svc.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
		return nil
	}
	for _, reservations := range result.Reservations {
		for _, instance := range reservations.Instances {
			//env, pillar, app, name := "", "", "", ""
			name := ""
			for _, tag := range instance.Tags{
				if *tag.Key == "Name" {
					name = *tag.Value
					continue
				}
				//if *tag.Key == "App" {
				//	app = *tag.Value
				//	continue
				//}
				//if *tag.Key == "Pillar" {
				//	pillar = *tag.Value
				//}
				//if *tag.Key == "Env" {
				//	env = *tag.Value
				//}
			}
			instances = append(instances, model.Instance{
				Address: *instance.PrivateIpAddress,
				Name: name,
				State: *instance.State.Name,
			})
		}
	}
	return instances
}


//UpdateCmdbServer server list update
func UpdateCmdbServer() {
	log.Println("start update service_list tables")
	instances := readsAWSServer()
	for _, instance := range instances{
		insertServerList(instance)
	}
}


func insertServerList(instance model.Instance) {
	sql := fmt.Sprintf(`
		INSERT INTO server_list (server, name, app, pillar, status, uptime)
		VALUES ('%s', '%s', '%s', '%s', '%s', now() at time zone 'utc')
		ON CONFLICT (server) 
		DO UPDATE SET
		server = EXCLUDED.server,
		name = EXCLUDED.name,
		app = EXCLUDED.app,
		pillar = EXCLUDED.pillar,
		uptime = EXCLUDED.uptime;`, instance.Address, instance.Name)
	db := database.Db()
	defer db.Close()
	row, err := db.Query(sql)
	defer row.Close()
	if err != nil {
		log.Printf("insert server_list error: %s - %v", instance.Address, instance.State, err)
	}
}
