package line

import (
	"encoding/json"
	"line-bot/aws"
	"line-bot/search"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Messenger interface {
	Reply(replyToken string, message *linebot.TextMessage) error
	EventRouter(events []*linebot.Event)
	ParseRequest(r events.APIGatewayProxyRequest) ([]*linebot.Event, error)
}

type message struct {
	ChannelSecret string
	ChannelToken  string
	Client        *linebot.Client
	AwsClient     *aws.Client
	search.Searcher
}

type User struct {
	ID     string `json:"id" dynamo:"user_line_id"`
	Status string `json:"status" dynamo:"status"`
	Genre  string `json:"genre" dynamo:"genre"`
	Place  string `json:"place" dynamo:"place"`
	Budget string `json:"budget" dynamo:"budget"`
}

func NewMessenger(secret, token, hotpepper_apikey, geocording_apikey string, _awsclient *aws.Client) (Messenger, error) {
	m := &message{
		ChannelSecret: secret,
		ChannelToken:  token,
	}

	m.Searcher = search.NewSearcher(hotpepper_apikey, geocording_apikey)
	m.AwsClient = _awsclient

	bot, err := linebot.New(
		m.ChannelSecret,
		m.ChannelToken,
	)
	if err != nil {
		return nil, err
	}
	m.Client = bot
	return m, nil
}

func (m *message) Reply(replyToken string, message *linebot.TextMessage) error {
	switch message.Text {
	case "りりこ":
		if _, err := m.Client.ReplyMessage(replyToken, linebot.NewTextMessage("がんばれ！！")).Do(); err != nil {
			return err
		}
	case "東京駅":
		shopAry, err := m.Searcher.Restaurant(message.Text, "5000", "フレンチ")
		if err != nil {
			return err
		}

		f := flexRestaurants(shopAry)
		if _, err := m.Client.ReplyMessage(replyToken, linebot.NewFlexMessage("検索結果", f)).Do(); err != nil {
			return err
		}

	default:
		if _, err := m.Client.ReplyMessage(replyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
			return err
		}
	}
	return nil
}

func (m *message) EventRouter(events []*linebot.Event) {
	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				err := m.statusCheck(event)
				if err != nil {
					log.Printf("DynamoDB Error: %v", err)
				}
				err = m.Reply(event.ReplyToken, message)
				if err != nil {
					log.Printf("Reply Error: %v", err)
				}
			}
		}
	}
}

func (m *message) ParseRequest(r events.APIGatewayProxyRequest) ([]*linebot.Event, error) {
	req := &struct {
		Events []*linebot.Event `json:"events"`
	}{}
	if err := json.Unmarshal([]byte(r.Body), req); err != nil {
		return nil, err
	}
	return req.Events, nil
}

func (m *message) statusCheck(event *linebot.Event) error {
	var user User
	err := m.AwsClient.IsLineUser(event.Source.UserID, user)
	if err != nil {
		return err
	}
	if user.ID != "" {
		user = User{ID: event.Source.UserID, Status: "WaitGenre"}
		err = m.AwsClient.SetLineUser(user)
		if err != nil {
			return err
		}
	}
	return nil
}
