package line

import (
	"encoding/json"
	"fmt"
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
		loc, err := m.Searcher.Place(message.Text)
		if err != nil {
			return err
		}
		responseLoc := linebot.NewLocationMessage(loc.Name, loc.Address, loc.Lat, loc.Lng)
		if _, err := m.Client.ReplyMessage(replyToken, responseLoc).Do(); err != nil {
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
				fmt.Println("-----debug-----")
				fmt.Printf("ReplyToken: %s\n", event.ReplyToken)
				fmt.Printf("UserID: %s\n", event.Source.UserID)
				fmt.Printf("GroupID: %s\n", event.Source.GroupID)
				fmt.Printf("RoomID: %s\n", event.Source.RoomID)
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
