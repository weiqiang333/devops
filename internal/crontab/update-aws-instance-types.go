package crontab

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/weiqiang333/devops/internal/awscli"
	"github.com/weiqiang333/devops/internal/database"
	"github.com/weiqiang333/devops/internal/model"
)


//UpdateAwsInstanceTypes 更新实例类型信息
func UpdateAwsInstanceTypes()  {
	log.Println("Start update aws_instance_types tables")
	svc := awscli.AwsCli()
	var instanceTypes = []model.TableAwsInstanceTypes{}
	instanceTypes, err := getAwsInstanceTypes(svc, "", instanceTypes)
	if err != nil {
		log.Printf("update aws_instance_types tables fail: %v", err)
	}
	err = insertInstanceTypes(instanceTypes)
	if err != nil {
		log.Printf("update aws_instance_types tables fail: %v", err)
	}
	log.Println("Success update aws_instance_types tables")
}

func getAwsInstanceTypes(svc *ec2.EC2, nextToken string, instanceTypes []model.TableAwsInstanceTypes) (
	[]model.TableAwsInstanceTypes, error) {
	input := &ec2.DescribeInstanceTypesInput{}
	if nextToken != "" {
		input = &ec2.DescribeInstanceTypesInput{
			NextToken: aws.String(nextToken),
		}
	}
	result, err := svc.DescribeInstanceTypes(input)
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
		return instanceTypes, fmt.Errorf("DescribeInstanceTypes fail: %v", err)
	}
	for _, types := range result.InstanceTypes {
		instanceTypes = append(instanceTypes, model.TableAwsInstanceTypes{
			InstanceType: *types.InstanceType,
			Vcpus: *types.VCpuInfo.DefaultVCpus,
			Memory: *types.MemoryInfo.SizeInMiB / 1024,
		})
	}
	if result.NextToken == nil {
		return instanceTypes, nil
	} else {
		return  getAwsInstanceTypes(svc, *result.NextToken, instanceTypes)
	}
}

func insertInstanceTypes(instanceTypes []model.TableAwsInstanceTypes) error {
	for i, instanceType := range instanceTypes {
		sql := fmt.Sprintf(`
			INSERT INTO aws_instance_types (instance_type, vcpus, memory, update_at)
			VALUES ('%s', %v, %v, now() at time zone 'utc')
			ON CONFLICT (instance_type) 
			DO UPDATE SET
			instance_type = EXCLUDED.instance_type,
			vcpus = EXCLUDED.vcpus,
			memory = EXCLUDED.memory,
			update_at = EXCLUDED.update_at;`, instanceType.InstanceType, instanceType.Vcpus, instanceType.Memory)
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
