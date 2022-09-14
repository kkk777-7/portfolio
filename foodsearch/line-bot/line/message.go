package line

import (
	"line-bot/search"
	"strconv"

	"github.com/line/line-bot-sdk-go/linebot"
)

func searchedRestaurants(shops []search.Shop) *linebot.CarouselContainer {
	return restaurants(shops, "お気に入り登録")
}

func favoriteRestaurants(shops []search.Shop) *linebot.CarouselContainer {
	return restaurants(shops, "お気に入り削除")
}

func restaurants(shops []search.Shop, footerMessage string) *linebot.CarouselContainer {
	var bcs []*linebot.BubbleContainer
	for _, shop := range shops {
		bc := linebot.BubbleContainer{
			Type:   linebot.FlexContainerTypeBubble,
			Hero:   setHero(shop),
			Body:   setBody(shop),
			Footer: setFooter(footerMessage, shop),
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

func setFooter(buttonmessage string, shop search.Shop) *linebot.BoxComponent {
	return &linebot.BoxComponent{
		Type:    linebot.FlexComponentTypeBox,
		Layout:  linebot.FlexBoxLayoutTypeVertical,
		Spacing: linebot.FlexComponentSpacingTypeXs,
		Contents: []linebot.FlexComponent{
			setUriButton("詳しく見る", shop.Url),
			setUriButton("地図を確認する", "https://www.google.com/maps"+"?q="+strconv.FormatFloat(shop.Lat, 'f', -1, 64)+","+strconv.FormatFloat(shop.Lng, 'f', -1, 64)),
			setUriButton("クーポンを確認", shop.Coupon),
			setPostButton(buttonmessage, shop),
		},
	}
}

func setUriButton(label string, uri string) *linebot.ButtonComponent {
	return &linebot.ButtonComponent{
		Type:   linebot.FlexComponentTypeButton,
		Style:  linebot.FlexButtonStyleTypeLink,
		Height: linebot.FlexButtonHeightTypeSm,
		Action: linebot.NewURIAction(label, uri),
	}
}

func setPostButton(label string, shop search.Shop) *linebot.ButtonComponent {
	return &linebot.ButtonComponent{
		Type:   linebot.FlexComponentTypeButton,
		Style:  linebot.FlexButtonStyleTypeLink,
		Height: linebot.FlexButtonHeightTypeSm,
		Action: linebot.NewPostbackAction(label, label+":"+shop.ShopId, "", ""),
	}
}

func quickReplyButton(labels []string, values []string) *linebot.QuickReplyItems {
	var buttons []*linebot.QuickReplyButton
	for i := 0; i < len(labels); i++ {
		button := linebot.QuickReplyButton{
			Action: &linebot.MessageAction{
				Label: labels[i],
				Text:  values[i],
			},
		}
		buttons = append(buttons, &button)
	}
	return &linebot.QuickReplyItems{Items: buttons}
}
