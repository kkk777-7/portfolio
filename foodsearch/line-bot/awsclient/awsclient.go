package awsclient

import (
	"line-bot/search"
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
	UserId string `json:"id"`
	Status string `json:"status"`
	Genre  string `json:"genre"`
	Place  string `json:"place"`
	Budget string `json:"budget"`
}

var table dynamo.Table

func NewClient() *Client {
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
	return client
}

func (c *Client) SsmGetParameter(key string) (string, error) {
	params := &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	}
	res, err := c.ssmsvc.GetParameter(params)
	if err != nil {
		return "", err
	}
	return *res.Parameter.Value, nil
}

func (c *Client) IsLineUser(tablename, id string, result *User) error {
	table = c.dynamosvc.Table(tablename)
	err := table.Get("UserId", id).One(result)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SetLineUser(tablename string, user *User) error {
	table = c.dynamosvc.Table(tablename)
	err := table.Put(user).Run()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateLineUser(tablename string, user *User, key, value string) error {
	table = c.dynamosvc.Table(tablename)
	err := table.Update("UserId", user.UserId).Set(key, value).Value(user)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SetShop(tablename string, shop *search.Shop) error {
	table = c.dynamosvc.Table(tablename)
	err := table.Put(shop).Run()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetShop(tablename, userid string, shops *[]search.Shop) error {
	table = c.dynamosvc.Table(tablename)
	err := table.Get("UserId", userid).All(shops)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteShop(tablename, userid, shopid string) error {
	table = c.dynamosvc.Table(tablename)
	err := table.Delete("UserId", userid).Range("ShopId", shopid).Run()
	if err != nil {
		return err
	}
	return nil
}
