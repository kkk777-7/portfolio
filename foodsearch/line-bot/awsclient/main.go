package awsclient

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/guregu/dynamo"
)

type Client struct {
	ssmsvc    *ssm.SSM
	dynamosvc *dynamo.DB
}

var table dynamo.Table

func NewClient() *Client {
	awsprofile := os.Getenv("AWS_PROFILE")

	client := new(Client)
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: awsprofile,
	}))
	client.ssmsvc = ssm.New(sess)

	if awsprofile == "local" {
		cfg := aws.NewConfig()
		cfg.Endpoint = aws.String(os.Getenv("DYNAMO_ENDPOINT"))
		client.dynamosvc = dynamo.New(sess, cfg)
	} else {
		client.dynamosvc = dynamo.New(sess)
	}
	table = client.dynamosvc.Table("users")
	return client
}

func (a *Client) SsmGetParameter(key string) (string, error) {
	params := &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	}
	res, err := a.ssmsvc.GetParameter(params)
	if err != nil {
		return "", err
	}
	return *res.Parameter.Value, nil
}
