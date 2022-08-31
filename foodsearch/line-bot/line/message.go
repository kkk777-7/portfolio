package line

import (
	"line-bot/search"
	"strconv"

	"github.com/line/line-bot-sdk-go/linebot"
)

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
			setButton("地図を確認する", "https://www.google.com/maps"+"?q="+strconv.FormatFloat(shop.Lat, 'f', -1, 64)+","+strconv.FormatFloat(shop.Lng, 'f', -1, 64)),
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
