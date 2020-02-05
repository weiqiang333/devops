package crontab

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/weiqiang333/devops/internal/database"
	"log"

	"github.com/weiqiang333/devops/internal/awscli"
	"github.com/weiqiang333/devops/internal/model"
)

//UpdateAwsInstanceTypes 更新卷信息
func UpdateAwsVolumes()  {
	log.Println("Start update aws_volumes tables")
	svc := awscli.AwsCli()
	var volumes = []model.TableAwsVolumes{}
	volumes, err := getAwsVolumes(svc, "", volumes)
	if err != nil {
		log.Printf("update aws_volumes tables fail: %v", err)
	}
	err = insertInstanceTypes(volumes)
	if err != nil {
		log.Printf("update aws_volumes tables fail: %v", err)
	}
	log.Println("Success update aws_volumes tables")
}

func getAwsVolumes(svc *ec2.EC2, nextToken string, volumes []model.TableAwsVolumes) ([]model.TableAwsVolumes, error) {
	input := &ec2.DescribeVolumesInput{}
	if nextToken != "" {
		input = &ec2.DescribeVolumesInput{
			NextToken: aws.String(nextToken),
		}
	}
	result, err := svc.DescribeVolumes(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return volumes, fmt.Errorf("getAwsVolumes fail: %v", err)
	}

	for _, volume := range result.Volumes {
		env, pillar, app, name := "", "", "", ""
		for _, tag := range volume.Tags{
			if *tag.Key == "Name" {
				name = *tag.Value
			}
			if *tag.Key == "App" {
				app = *tag.Value
			}
			if *tag.Key == "Pillar" {
				pillar = *tag.Value
			}
			if *tag.Key == "Env" {
				env = *tag.Value
			}
		}
		volumes = append(volumes, model.TableAwsVolumes{
			DiskId: *volume.VolumeId,
			CreateAt: *volume.CreateTime,
			Size: *volume.Size,
			Iops: *volume.Iops,
			State: *volume.State,
			Type: *volume.VolumeType,
			TagName: name,
			TagApp: app,
			TagEnv: env,
			TagPillar: pillar,
		})
	}
	if result.NextToken == nil {
		return volumes, nil
	} else {
		return  getAwsVolumes(svc, *result.NextToken, volumes)
	}
}

func insertAwsVolumes(volumes []model.TableAwsVolumes) error {
	for i, volume := range volumes {
		sql := fmt.Sprintf(`
			INSERT INTO aws_volumes (disk_id, create_at, size, iops, state, type, tag_name, tag_app, tag_env, tag_pillar, update_at)
			VALUES ('%s', %v, %v, %v, '%s', '%s', now() at time zone 'utc')
			ON CONFLICT (instance_type)
			DO UPDATE SET
			instance_type = EXCLUDED.instance_type,
			vcpus = EXCLUDED.vcpus,
			memory = EXCLUDED.memory,
			update_at = EXCLUDED.update_at;`, volume.DiskId, volume.CreateAt, volume.Size, volume.Iops, volume.State, volume.Type)
		db := database.Db()
		defer db.Close()
		row, err := db.Query(sql)
		defer row.Close()
		if err != nil {
			log.Printf("insert aws_instance_types error: %s - %v", instanceType.InstanceType, err)
			return fmt.Errorf("insert aws_instance_types error: %s - %v", instanceType.InstanceType, err)
		}
		log.Printf("insert aws_instance_types %v %s", i, instanceType.InstanceType)
	}
	return nil
}
