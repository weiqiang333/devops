package awscli

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)


// AwsCli conn aws api
func AwsCli() *ec2.EC2 {
	sess := session.New()
	svc := ec2.New(sess, &aws.Config{
		Region:	aws.String(endpoints.CnNorth1RegionID),
	})
	return svc
}