package search

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/itchyny/gojq"
	"github.com/pkg/errors"
)

const (
	HOTPEPPER_APIENDPOINT  = "https://webservice.recruit.co.jp/hotpepper/gourmet/v1/"
	GEOCORDING_APIENDPOINT = "https://maps.googleapis.com/maps/api/geocode/json"

	HOTPEPPER_JQ_QUERY  = ".results.shop[] | { shopid: .id, name: .name, access: .mobile_access, budget: .budget.average, url: .urls.pc, photo: .photo.mobile.l, lat: .lat, lng: .lng, coupon: .coupon_urls.sp, genre: .genre.name }"
	GEOCORDING_JQ_QUERY = ".results[] | { address: .formatted_address, lat: .geometry.location.lat, lng: .geometry.location.lng}"
)

type Searcher interface {
	Restaurant(place, budget, genre string) ([]Shop, error)
	RestaurantById(id string) (*Shop, error)
	Place(place string) (*Location, error)
}

type search struct {
	hotpepperApiUrl  string
	hotpepperApiKey  string
	geocordingApiUrl string
	geocordingApiKey string
}

type Shop struct {
	UserId string  `json:"id"`
	ShopId string  `json:"shopid"`
	Name   string  `json:"name"`
	Access string  `json:"access"`
	Budget string  `json:"budget"`
	Url    string  `json:"url"`
	Photo  string  `json:"photo"`
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	Coupon string  `json:"coupon"`
	Genre  string  `json:"genre"`
}

type Location struct {
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Address string  `json:"address"`
}

// Constructor of Searcher from API
func NewSearcher(_hotpepperApiUrl, _hotpepperApiKey, _geocordingApiUrl, _geocordingApiKey string) Searcher {
	hotpepperUrl := HOTPEPPER_APIENDPOINT
	geocordingUrl := GEOCORDING_APIENDPOINT
	if _hotpepperApiUrl != "" {
		hotpepperUrl = _hotpepperApiUrl
	}
	if _geocordingApiUrl != "" {
		geocordingUrl = _geocordingApiUrl
	}

	return &search{
		hotpepperApiUrl:  hotpepperUrl,
		hotpepperApiKey:  _hotpepperApiKey,
		geocordingApiUrl: geocordingUrl,
		geocordingApiKey: _geocordingApiKey,
	}
}

// Search for the  restaurant from hotpepper's API
func (s *search) Restaurant(place, budget, genre string) ([]Shop, error) {
	loc, err := s.Place(place)
	if err != nil {
		return nil, err
	}
	budgetcode, err := convertValueToCode(budget)
	if err != nil {
		return nil, err
	}
	genrecode := convertGenreToCode(genre)

	params := url.Values{}
	params.Add("key", s.hotpepperApiKey)
	params.Add("lat", strconv.FormatFloat(loc.Lat, 'f', -1, 64))
	params.Add("lng", strconv.FormatFloat(loc.Lng, 'f', -1, 64))
	params.Add("range", "3")
	params.Add("budget", budgetcode)
	params.Add("genre", genrecode)
	params.Add("count", "5")
	params.Add("format", "json")

	resp, err := http.Get(HOTPEPPER_APIENDPOINT + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResult interface{}
	if err := json.Unmarshal(body, &apiResult); err != nil {
		return nil, err
	}

	shopAry := []Shop{}
	query, err := gojq.Parse(HOTPEPPER_JQ_QUERY)
	if err != nil {
		return nil, err
	}

	iter := query.Run(apiResult)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		var shop Shop
		if err := json.Unmarshal(jsonBytes, &shop); err != nil {
			return nil, err
		}
		shopAry = append(shopAry, shop)
	}
	return shopAry, nil
}

// Search for the  restaurant By Id from hotpepper's API
func (s *search) RestaurantById(id string) (*Shop, error) {
	params := url.Values{}
	params.Add("key", s.hotpepperApiKey)
	params.Add("id", id)
	params.Add("format", "json")

	resp, err := http.Get(HOTPEPPER_APIENDPOINT + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResult interface{}
	if err := json.Unmarshal(body, &apiResult); err != nil {
		return nil, err
	}

	var shop Shop
	query, err := gojq.Parse(HOTPEPPER_JQ_QUERY)
	if err != nil {
		return nil, err
	}

	iter := query.Run(apiResult)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(jsonBytes, &shop); err != nil {
			return nil, err
		}
	}
	return &shop, nil
}

// Search the latitude and longitude of a point from geocording's API
func (s *search) Place(place string) (*Location, error) {
	params := url.Values{}
	params.Add("address", place)
	params.Add("key", s.geocordingApiKey)

	resp, err := http.Get(s.geocordingApiUrl + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResult interface{}
	if err := json.Unmarshal(body, &apiResult); err != nil {
		return nil, err
	}

	var loc Location
	query, err := gojq.Parse(GEOCORDING_JQ_QUERY)
	if err != nil {
		return nil, err
	}

	iter := query.Run(apiResult)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		loc.Name = place
		if err := json.Unmarshal(jsonBytes, &loc); err != nil {
			return nil, err
		}
	}
	return &loc, nil
}

// Convert from budget to API code
func convertValueToCode(budget string) (string, error) {
	code, err := strconv.Atoi(budget)
	if err != nil {
		return "", err
	}
	switch {
	case 0 < code && code <= 500:
		return "B009", nil
	case 500 < code && code <= 1000:
		return "B010", nil
	case 1000 < code && code <= 1500:
		return "B011", nil
	case 1500 < code && code <= 2000:
		return "B001", nil
	case 2000 < code && code <= 3000:
		return "B002", nil
	case 3000 < code && code <= 4000:
		return "B003", nil
	case 4000 < code && code <= 5000:
		return "B008", nil
	case 5000 < code && code <= 7000:
		return "B004", nil
	case 7000 < code && code <= 10000:
		return "B005", nil
	case 10000 < code && code <= 15000:
		return "B006", nil
	case 15000 < code && code <= 20000:
		return "B012", nil
	case 20000 < code && code <= 30000:
		return "B013", nil
	case 30000 < code:
		return "B014", nil
	}
	return "", errors.New("Invalid value.")
}

//　Convert from genre name to API code
func convertGenreToCode(genre string) string {
	switch {
	case strings.Contains(genre, "居酒屋"):
		return "G001"
	case strings.Contains(genre, "ダイニングバー"):
		return "G002"
	case strings.Contains(genre, "創作料理"):
		return "G003"
	case strings.Contains(genre, "和食"):
		return "G004"
	case strings.Contains(genre, "洋食"):
		return "G005"
	case strings.Contains(genre, "イタリアン") || strings.Contains(genre, "フレンチ"):
		return "G006"
	case strings.Contains(genre, "中華"):
		return "G007"
	case strings.Contains(genre, "焼肉"):
		return "G008"
	case strings.Contains(genre, "エスニック料理"):
		return "G009"
	case strings.Contains(genre, "各国料理"):
		return "G010"
	case strings.Contains(genre, "カラオケ"):
		return "G011"
	case strings.Contains(genre, "バー"):
		return "G012"
	case strings.Contains(genre, "ラーメン"):
		return "G013"
	case strings.Contains(genre, "カフェ"):
		return "G014"
	case strings.Contains(genre, "お好み焼き"):
		return "G016"
	case strings.Contains(genre, "韓国料理"):
		return "G017"
	}
	return "G015"
}
