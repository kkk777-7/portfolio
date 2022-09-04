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
	UserId string `json:"id" dynamo:"UserId"`
	Status string `json:"status" dynamo:"Status"`
	Genre  string `json:"genre" dynamo:"Genre"`
	Place  string `json:"place" dynamo:"Place"`
	Budget string `json:"budget" dynamo:"Budget"`
}

var table dynamo.Table

func NewClient(tablename string) *Client {
	client := new(Client)

	disableSsl := false
	awsRegion := os.Getenv("AWSREGION")
	dynamoDbEndpoint := os.Getenv("DYNAMOENDPOINT")

	if len(awsRegion) == 0 {
		awsRegion = "ap-northeast-1"
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)}))
	client.ssmsvc = ssm.New(sess)

	if len(dynamoDbEndpoint) != 0 {
		disableSsl = true
		client.dynamosvc = dynamo.New(session.New(), &aws.Config{
			Region:     aws.String(awsRegion),
			Endpoint:   aws.String(dynamoDbEndpoint),
			DisableSSL: aws.Bool(disableSsl),
		})
	} else {
		client.dynamosvc = dynamo.New(session.New(), &aws.Config{
			Region:     aws.String(awsRegion),
			DisableSSL: aws.Bool(disableSsl),
		})
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

func (a *Client) IsLineUser(id string, result *User) error {
	err := table.Get("UserId", id).One(result)
	if err != nil {
		return err
	}
	return nil
}

func (a *Client) SetLineUser(user *User) error {
	err := table.Put(user).Run()
	if err != nil {
		return err
	}
	return nil
}

func (a *Client) UpdateLineUser(user *User, key, value string) error {
	err := table.Update("UserId", user.UserId).Set(key, value).Value(user)
	if err != nil {
		return err
	}
	return nil
}
