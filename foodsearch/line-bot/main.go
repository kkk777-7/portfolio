package main

import (
	"encoding/json"
	"os"

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

func main() {
	lambda.Start(Handler)
}

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r := &Line{}
	err := r.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_TOKEN"))
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
