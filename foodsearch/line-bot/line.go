package main

import (
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

type Line struct {
	ChannelSecret string
	ChannelToken  string
	Client        *linebot.Client
}

func (r *Line) New(secret, token string) error {
	r.ChannelSecret = secret
	r.ChannelToken = token

	bot, err := linebot.New(
		r.ChannelSecret,
		r.ChannelToken,
	)
	if err != nil {
		return err
	}
	r.Client = bot
	return nil
}

func (r *Line) Reply(replyToken string, message *linebot.TextMessage) error {
	if _, err := r.Client.ReplyMessage(replyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
		return err
	}
	return nil
}

func (r *Line) EventRouter(events []*linebot.Event) {
	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				err := r.Reply(event.ReplyToken, message)
				log.Printf("Reply Error: %v", err)
			}
		}
	}
}
