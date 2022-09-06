package line

import (
	"encoding/json"
	"fmt"
	"line-bot/awsclient"
	"line-bot/search"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Messenger interface {
	EventRouter(events []*linebot.Event)
	Reply(event *linebot.Event, message *linebot.TextMessage) error
	ParseRequest(r events.APIGatewayProxyRequest) ([]*linebot.Event, error)
}

type message struct {
	ChannelSecret string
	ChannelToken  string
	Client        *linebot.Client
	AwsClient     *awsclient.Client
	search.Searcher
}

var user awsclient.User
var genres []string = []string{"居酒屋", "和食", "洋食", "イタリアン", "フレンチ", "中華", "焼肉", "カラオケ", "バー", "ラーメン", "カフェ", "その他"}
var budgets []string = []string{"1000", "1500", "2000", "3000", "4000", "5000", "7000", "10000", "15000", "20000", "30000", "30001"}
var budgetLabels []string = []string{"500円以上1000円未満", "1000円以上1500円未満", "1500円以上2000円未満", "2000円以上3000円未満", "3000円以上4000円未満", "4000円以上5000円未満", "5000円以上7000円未満", "7000円以上10000円未満", "10000円以上15000円未満", "15000円以上20000円未満", "20000円以上30000円未満", "30000円以上"}

func NewMessenger(secret, token, hotpepper_apikey, geocording_apikey string, _awsclient *awsclient.Client) (Messenger, error) {
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

func (m *message) EventRouter(events []*linebot.Event) {
	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				err := m.Reply(event, message)
				if err != nil {
					log.Printf("Reply Error: %v", err)
				}
			}
		case linebot.EventTypePostback:
			err := m.Register(event)
			if err != nil {
				log.Printf("Register Error: %v", err)
			}
		}
	}
}

func (m *message) Reply(event *linebot.Event, message *linebot.TextMessage) error {
	replyToken := event.ReplyToken

	if message.Text == "お気に入り表示" {
		var shopAry []search.Shop
		err := m.AwsClient.GetShop("shops", event.Source.UserID, &shopAry)
		if err != nil {
			return err
		}
		f := flexRestaurants(shopAry)
		if _, err := m.Client.ReplyMessage(replyToken, linebot.NewFlexMessage("検索結果", f)).Do(); err != nil {
			return err
		}
		return nil
	}

	if err := m.statusCheck(event); err != nil {
		return err
	}
	if user.Status == "WaitSearch" && message.Text != "検索" {
		if _, err := m.Client.ReplyMessage(replyToken, linebot.NewTextMessage("お店を検索したい場合は、「検索」と入力してね！")).Do(); err != nil {
			return err
		}
	} else {
		switch user.Status {
		case "WaitSearch":
			if _, err := m.Client.ReplyMessage(replyToken, linebot.NewTextMessage("どこで食べる？\n住所や地名で検索してね！")).Do(); err != nil {
				return err
			}
			err := m.AwsClient.UpdateLineUser("users", &user, "Status", "WaitPlace")
			if err != nil {
				return err
			}
		case "WaitPlace":
			replyMessage := linebot.NewTextMessage("どんなジャンル？\n下から選ぶ または 入力してね！")
			buttons := quickReplyButton(genres, genres)
			replyMessage.WithQuickReplies(buttons)
			if _, err := m.Client.ReplyMessage(replyToken, replyMessage).Do(); err != nil {
				return err
			}
			err := m.AwsClient.UpdateLineUser("users", &user, "Place", message.Text)
			err = m.AwsClient.UpdateLineUser("users", &user, "Status", "WaitGenre")
			if err != nil {
				return err
			}
		case "WaitGenre":
			replyMessage := linebot.NewTextMessage("どのくらいの予算？\n下から選ぶ または 入力してね！")
			buttons := quickReplyButton(budgetLabels, budgets)
			replyMessage.WithQuickReplies(buttons)
			if _, err := m.Client.ReplyMessage(replyToken, replyMessage).Do(); err != nil {
				return err
			}
			err := m.AwsClient.UpdateLineUser("users", &user, "Genre", message.Text)
			err = m.AwsClient.UpdateLineUser("users", &user, "Status", "WaitBudget")
			if err != nil {
				return err
			}
		case "WaitBudget":
			err := m.AwsClient.UpdateLineUser("users", &user, "Budget", message.Text)
			err = m.AwsClient.UpdateLineUser("users", &user, "Status", "Searching")
			if err != nil {
				return err
			}

			shopAry, err := m.Searcher.Restaurant(user.Place, user.Budget, user.Genre)
			if err != nil {
				return err
			}

			for i := 0; i < len(shopAry); i++ {
				if shopAry[i].Access == "" {
					shopAry[i].Access = "-"
				}
				if shopAry[i].Budget == "" {
					shopAry[i].Budget = "-"
				}
			}

			if len(shopAry) == 0 {
				if _, err := m.Client.ReplyMessage(replyToken, linebot.NewTextMessage("お店が見つかりませんでした...\n条件を見直してね！")).Do(); err != nil {
					return err
				}
			} else {
				f := flexRestaurants(shopAry)
				if _, err := m.Client.ReplyMessage(replyToken, linebot.NewFlexMessage("検索結果", f)).Do(); err != nil {
					return err
				}
			}
			if err := m.statusReset(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *message) Register(event *linebot.Event) error {
	replyToken := event.ReplyToken

	shop_id := event.Postback.Data
	shop, err := m.Searcher.RestaurantById(shop_id)
	if err != nil {
		return err
	}

	if shop.Access == "" {
		shop.Access = "-"
	}
	if shop.Budget == "" {
		shop.Budget = "-"
	}

	shop.UserId = event.Source.UserID
	fmt.Printf("%+v\n", shop)
	err = m.AwsClient.SetShop("shops", &shop)
	if err != nil {
		return err
	}

	replyMessage := fmt.Sprintf("%sをお気に入り登録したよ！", shop.Name)
	if _, err := m.Client.ReplyMessage(replyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
		return err
	}
	return nil
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
	err := m.AwsClient.IsLineUser("users", event.Source.UserID, &user)
	if err != nil && user.UserId != "" {
		return err
	}
	if user.UserId == "" {
		user = awsclient.User{UserId: event.Source.UserID, Status: "WaitSearch"}
		err = m.AwsClient.SetLineUser("users", &user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *message) statusReset() error {
	err := m.AwsClient.UpdateLineUser("users", &user, "Status", "WaitSearch")
	if err != nil {
		return err
	}
	return nil
}
