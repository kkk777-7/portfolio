package line

import (
	"encoding/json"
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
	search.Searcher
}

func NewMessenger(secret, token, hotpepper_apikey, geocording_apikey string) (Messenger, error) {
	m := &message{
		ChannelSecret: secret,
		ChannelToken:  token,
	}

	m.Searcher = search.NewSearcher(hotpepper_apikey, geocording_apikey)

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
		for i := 0; i < len(shopAry); i++ {
			replyMessage := shopAry[i].Name + "\n" + shopAry[i].Address + "\n" + shopAry[i].Open + "\n" + shopAry[i].Url
			if _, err := m.Client.ReplyMessage(replyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
				return err
			}
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
				err := m.Reply(event.ReplyToken, message)
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
