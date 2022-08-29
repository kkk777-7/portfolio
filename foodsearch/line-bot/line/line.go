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

func flexRestaurants(shops []search.Shop) *linebot.CarouselContainer {
	var bcs []*linebot.BubbleContainer
	for _, shop := range shops {
		bc := linebot.BubbleContainer{
			Type:   linebot.FlexContainerTypeBubble,
			Hero:   setHero(shop),
			Body:   setBody(shop),
			Footer: setFooter(shop),
		}
		bcs = append(bcs, &bc)
	}
	return &linebot.CarouselContainer{
		Type:     linebot.FlexContainerTypeCarousel,
		Contents: bcs,
	}
}

func setHero(shop search.Shop) *linebot.ImageComponent {
	if shop.Photo == "" {
		return nil
	}
	return &linebot.ImageComponent{
		Type:        linebot.FlexComponentTypeImage,
		URL:         shop.Photo,
		Size:        linebot.FlexImageSizeTypeFull,
		AspectRatio: linebot.FlexImageAspectRatioType20to13,
		AspectMode:  linebot.FlexImageAspectModeTypeCover,
	}
}

func setBody(shop search.Shop) *linebot.BoxComponent {
	return &linebot.BoxComponent{
		Type:   linebot.FlexComponentTypeBox,
		Layout: linebot.FlexBoxLayoutTypeVertical,
		Contents: []linebot.FlexComponent{
			setRestaurantName(shop),
			setLocation(shop),
			setBudget(shop),
		},
	}
}

func setRestaurantName(shop search.Shop) *linebot.TextComponent {
	return &linebot.TextComponent{
		Type:   linebot.FlexComponentTypeText,
		Text:   shop.Name,
		Wrap:   true,
		Weight: linebot.FlexTextWeightTypeBold,
		Size:   linebot.FlexTextSizeTypeMd,
	}
}

func setLocation(shop search.Shop) *linebot.BoxComponent {
	return &linebot.BoxComponent{
		Type:    linebot.FlexComponentTypeBox,
		Layout:  linebot.FlexBoxLayoutTypeBaseline,
		Margin:  linebot.FlexComponentMarginTypeLg,
		Spacing: linebot.FlexComponentSpacingTypeSm,
		Contents: []linebot.FlexComponent{
			setSubtitle("エリア"),
			setDetail(shop.Access),
		},
	}
}

func setBudget(shop search.Shop) *linebot.BoxComponent {
	return &linebot.BoxComponent{
		Type:    linebot.FlexComponentTypeBox,
		Layout:  linebot.FlexBoxLayoutTypeBaseline,
		Margin:  linebot.FlexComponentMarginTypeLg,
		Spacing: linebot.FlexComponentSpacingTypeSm,
		Contents: []linebot.FlexComponent{
			setSubtitle("予算"),
			setDetail(shop.Budget),
		},
	}
}

func setSubtitle(t string) *linebot.TextComponent {
	return &linebot.TextComponent{
		Type:  linebot.FlexComponentTypeText,
		Text:  t,
		Color: "#aaaaaa",
		Size:  linebot.FlexTextSizeTypeXs,
		Flex:  linebot.IntPtr(4),
	}
}

func setDetail(t string) *linebot.TextComponent {
	return &linebot.TextComponent{
		Type:  linebot.FlexComponentTypeText,
		Text:  t,
		Wrap:  true,
		Color: "#666666",
		Size:  linebot.FlexTextSizeTypeXs,
		Flex:  linebot.IntPtr(12),
	}
}

func setFooter(shop search.Shop) *linebot.BoxComponent {
	return &linebot.BoxComponent{
		Type:    linebot.FlexComponentTypeBox,
		Layout:  linebot.FlexBoxLayoutTypeVertical,
		Spacing: linebot.FlexComponentSpacingTypeXs,
		Contents: []linebot.FlexComponent{
			setButton("詳しく見る", shop.Url),
			setButton("地図を確認する", "https://www.google.com/maps"+"?q="+shop.Lat+","+shop.Lng),
			setButton("クーポンを確認", shop.Coupon),
		},
	}
}

func setButton(label string, uri string) *linebot.ButtonComponent {
	return &linebot.ButtonComponent{
		Type:   linebot.FlexComponentTypeButton,
		Style:  linebot.FlexButtonStyleTypeLink,
		Height: linebot.FlexButtonHeightTypeSm,
		Action: linebot.NewURIAction(label, uri),
	}
}
