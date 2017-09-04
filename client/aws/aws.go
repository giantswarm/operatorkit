package aws

import (
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/giantswarm/microerror"
)

type Config struct {
	AccessKeyID     string
	AccessKeySecret string
	Region          string
	accountID       string
}

type Clients struct {
	AutoScaling    *autoscaling.AutoScaling
	CloudFormation *cloudformation.CloudFormation
	EC2            *ec2.EC2
	ELB            *elb.ELB
	IAM            *iam.IAM
	KMS            *kms.KMS
	Route53        *route53.Route53
	S3             *s3.S3
}

const (
	accountIDPosition = 4
	accountIDLength   = 12
)

func NewClients(config Config) Clients {
	awsCfg := &aws.Config{
		Credentials: credentials.NewStaticCredentials(config.AccessKeyID, config.AccessKeySecret, ""),
		Region:      aws.String(config.Region),
	}
	s := session.New(awsCfg)
	clients := Clients{
		AutoScaling:    autoscaling.New(s),
		CloudFormation: cloudformation.New(s),
		EC2:            ec2.New(s),
		ELB:            elb.New(s),
		IAM:            iam.New(s),
		KMS:            kms.New(s),
		Route53:        route53.New(s),
		S3:             s3.New(s),
	}

	return clients
}

func (c *Config) SetAccountID(iamClient *iam.IAM) error {
	resp, err := iamClient.GetUser(&iam.GetUserInput{})
	if err != nil {
		return microerror.Mask(err)
	}

	userArn := *resp.User.Arn
	accountID := strings.Split(userArn, ":")[accountIDPosition]

	if err := validateAccountID(accountID); err != nil {
		return microerror.Mask(err)
	}

	c.accountID = accountID

	return nil
}

func (c *Config) AccountID() string {
	return c.accountID
}

func validateAccountID(accountID string) error {
	r, _ := regexp.Compile("^[0-9]*$")

	switch {
	case accountID == "":
		return microerror.Mask(emptyAmazonAccountIDError)
	case len(accountID) != accountIDLength:
		return microerror.Mask(wrongAmazonAccountIDLengthError)
	case !r.MatchString(accountID):
		return microerror.Mask(malformedAmazonAccountIDError)
	}

	return nil
}
