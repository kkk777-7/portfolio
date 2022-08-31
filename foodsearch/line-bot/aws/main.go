package aws

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

func NewClient(tablename string) *Client {
	awsprofile := os.Getenv("AWSPROFILE")

	client := new(Client)
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: awsprofile,
	}))
	client.ssmsvc = ssm.New(sess)

	if awsprofile == "local" {
		cfg := aws.NewConfig()
		cfg.Endpoint = aws.String(os.Getenv("DYNAMOENDPOINT"))
		client.dynamosvc = dynamo.New(sess, cfg)
	} else {
		client.dynamosvc = dynamo.New(sess)
	}
	table = client.dynamosvc.Table(tablename)
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

func (a *Client) IsLineUser(id string, results interface{}) error {
	err := table.Get("user_line_id", id).All(&results)
	if err != nil {
		return err
	}
	return nil
}

func (a *Client) SetLineUser(user interface{}) error {
	err := table.Put(user).Run()
	if err != nil {
		return err
	}
	return nil
}
