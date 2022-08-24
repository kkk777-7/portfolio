package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/line/line-bot-sdk-go/linebot"
)

const (
	Successful = 200
	BadReq     = 400
	ErrSsm     = 500
	ErrReq     = 500
)

var CHANNEL_SECRET string
var CHANNEL_TOKEN string

func init() {
	err := setupParameters()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	lambda.Start(Handler)
}

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println(CHANNEL_SECRET)
	fmt.Println("------------")
	fmt.Println(CHANNEL_TOKEN)
	r := &Line{}
	err := r.New(CHANNEL_SECRET, CHANNEL_TOKEN)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: ErrSsm,
		}, nil
	}
	event, err := parseRequest(r.ChannelSecret, req)
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

	r.EventRouter(event)
	return events.APIGatewayProxyResponse{
		Body:       req.Body,
		StatusCode: Successful,
	}, nil
}

func parseRequest(channelSecret string, r events.APIGatewayProxyRequest) ([]*linebot.Event, error) {
	req := &struct {
		Events []*linebot.Event `json:"events"`
	}{}
	if err := json.Unmarshal([]byte(r.Body), req); err != nil {
		return nil, err
	}
	return req.Events, nil
}

func setupParameters() error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1")}))
	svc := ssm.New(sess)

	params := &ssm.GetParameterInput{
		Name:           aws.String("channel_secret_testbot"),
		WithDecryption: aws.Bool(true),
	}
	res, err := svc.GetParameter(params)
	if err != nil {
		return err
	}
	CHANNEL_SECRET = *res.Parameter.Value

	params = &ssm.GetParameterInput{
		Name:           aws.String("channel_token_testbot"),
		WithDecryption: aws.Bool(true),
	}
	res, err = svc.GetParameter(params)
	if err != nil {
		return err
	}
	CHANNEL_TOKEN = *res.Parameter.Value
	return nil
}
