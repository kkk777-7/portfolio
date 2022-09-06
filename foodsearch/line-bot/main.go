package main

import (
	"line-bot/awsclient"
	"line-bot/line"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
)

const (
	Successful = 200
	BadReq     = 400
	ErrSsm     = 500
	ErrReq     = 500
)

var awsClient *awsclient.Client
var lineClinet line.Messenger
var CHANNEL_SECRET string
var CHANNEL_TOKEN string
var HOTPEPPER_KEY string
var GOOGLE_KEY string

func init() {
	awsClient = awsclient.NewClient()
	err := setupParameters()
	if err != nil {
		log.Fatal(err)
	}
	lineClinet, err = line.NewMessenger(CHANNEL_SECRET, CHANNEL_TOKEN, HOTPEPPER_KEY, GOOGLE_KEY, awsClient)
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
	CHANNEL_SECRET, err = awsClient.SsmGetParameter("channel_secret_testbot")
	if err != nil {
		return err
	}
	CHANNEL_TOKEN, err = awsClient.SsmGetParameter("channel_token_testbot")
	if err != nil {
		return err
	}
	HOTPEPPER_KEY, err = awsClient.SsmGetParameter("hotpepper_key")
	if err != nil {
		return err
	}
	GOOGLE_KEY, err = awsClient.SsmGetParameter("google_key")
	if err != nil {
		return err
	}
	return nil
}
