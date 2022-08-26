package search

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/itchyny/gojq"
	"github.com/pkg/errors"
)

const (
	HOTPEPPER_APIENDPOINT  = "https://webservice.recruit.co.jp/hotpepper/gourmet/v1/"
	GEOCORDING_APIENDPOINT = "https://maps.googleapis.com/maps/api/geocode/json"

	GEOCORDING_JQ_QUERY = ".results[] | { address: .formatted_address, lat: .geometry.location.lat, lng: .geometry.location.lng}"
)

type Searcher interface {
	Place(place string) (*Location, error)
}

type search struct {
	hotpepper_apikey  string
	geocording_apikey string
}

type Location struct {
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Address string  `json:"address"`
}

func NewSearcher(_hotpepper_apikey, _geocording_apikey string) Searcher {
	return &search{
		hotpepper_apikey:  _hotpepper_apikey,
		geocording_apikey: _geocording_apikey,
	}
}

/*
func Restaurants(w string) {
	v := url.Values{}
	v.Add("range", "2")
	v.Add("order", "4")
	v.Add("format", "json")
	fmt.Println(v.Encode())
}
*/

func (s *search) Place(place string) (*Location, error) {
	params := url.Values{}
	params.Add("address", place)
	params.Add("key", s.geocording_apikey)

	resp, err := http.Get(GEOCORDING_APIENDPOINT + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResult interface{}
	if err := json.Unmarshal(body, &apiResult); err != nil {
		return nil, err
	}

	query, err := gojq.Parse(GEOCORDING_JQ_QUERY)
	if err != nil {
		return nil, err
	}

	iter := query.Run(apiResult)
	v, _ := iter.Next()

	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	loc := &Location{Name: place}
	if err := json.Unmarshal(jsonBytes, loc); err != nil {
		return nil, err
	}

	return loc, nil
}

func convertValueToCode(value string) (string, error) {
	chk, err := strconv.Atoi(value)
	if err != nil {
		return "", err
	}
	switch {
	case 0 < chk && chk <= 500:
		return "B009", nil
	case 500 < chk && chk <= 1000:
		return "B010", nil
	case 1000 < chk && chk <= 1500:
		return "B011", nil
	case 1500 < chk && chk <= 2000:
		return "B001", nil
	case 2000 < chk && chk <= 3000:
		return "B002", nil
	case 3000 < chk && chk <= 4000:
		return "B003", nil
	case 4000 < chk && chk <= 5000:
		return "B008", nil
	case 5000 < chk && chk <= 7000:
		return "B004", nil
	case 7000 < chk && chk <= 10000:
		return "B005", nil
	case 10000 < chk && chk <= 15000:
		return "B006", nil
	case 15000 < chk && chk <= 20000:
		return "B012", nil
	case 20000 < chk && chk <= 30000:
		return "B013", nil
	case 30000 < chk:
		return "B014", nil
	}
	return "", errors.New("Invalid value.")
}
