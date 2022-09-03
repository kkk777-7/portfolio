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

type User struct {
	ID     string `json:"id" dynamo:"user_line_id"`
	Status string `json:"status" dynamo:"status"`
	Genre  string `json:"genre" dynamo:"genre"`
	Place  string `json:"place" dynamo:"place"`
	Budget string `json:"budget" dynamo:"budget"`
}

var table dynamo.Table

func NewClient(tablename string) *Client {
	disableSsl := false
	dynamoDbRegion := os.Getenv("AWSREGION")
	dynamoDbEndpoint := os.Getenv("DYNAMOENDPOINT")
	if len(dynamoDbEndpoint) != 0 {
		disableSsl = true
	}
	if len(dynamoDbRegion) == 0 {
		dynamoDbRegion = "ap-northeast-1"
	}

	client := new(Client)
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1")}))
	client.ssmsvc = ssm.New(sess)

	client.dynamosvc = dynamo.New(session.New(), &aws.Config{
		Region:     aws.String(dynamoDbRegion),
		Endpoint:   aws.String(dynamoDbEndpoint),
		DisableSSL: aws.Bool(disableSsl),
	})

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

func (a *Client) IsLineUser(id string, result *User) error {
	err := table.Get("user_line_id", id).One(result)
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
