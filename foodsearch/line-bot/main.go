package main

import (
	"line-bot/line"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/line/line-bot-sdk-go/linebot"
)

type AWSClient struct {
	ssmsvc *ssm.SSM
}

const (
	Successful = 200
	BadReq     = 400
	ErrSsm     = 500
	ErrReq     = 500
)

var awsClient *AWSClient
var lineClinet line.Messenger
var CHANNEL_SECRET string
var CHANNEL_TOKEN string
var HOTPEPPER_KEY string
var GOOGLE_KEY string

func init() {
	awsClient = NewAWSClient()
	err := setupParameters()
	if err != nil {
		log.Fatal(err)
	}
	lineClinet, err = line.NewMessenger(CHANNEL_SECRET, CHANNEL_TOKEN, HOTPEPPER_KEY, GOOGLE_KEY)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	lambda.Start(Handler)
}

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	event, err := lineClinet.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			return events.APIGatewayProxyResponse{
				Body:       err.Error(),
				StatusCode: BadReq,
			}, nil
		} else {
			return events.APIGatewayProxyResponse{
				Body:       err.Error(),
				StatusCode: ErrReq,
			}, nil
		}
	}

	lineClinet.EventRouter(event)
	return events.APIGatewayProxyResponse{
		StatusCode: Successful,
	}, nil
}

func setupParameters() error {
	var err error
	CHANNEL_SECRET, err = awsClient.ssmGetParameter("channel_secret_testbot")
	if err != nil {
		return err
	}
	CHANNEL_TOKEN, err = awsClient.ssmGetParameter("channel_token_testbot")
	if err != nil {
		return err
	}
	HOTPEPPER_KEY, err = awsClient.ssmGetParameter("hotpepper_key")
	if err != nil {
		return err
	}
	GOOGLE_KEY, err = awsClient.ssmGetParameter("google_key")
	if err != nil {
		return err
	}
	return nil
}

func (a *AWSClient) ssmGetParameter(key string) (string, error) {
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

func NewAWSClient() *AWSClient {
	client := new(AWSClient)
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1")}))
	client.ssmsvc = ssm.New(sess)

	return client
}
