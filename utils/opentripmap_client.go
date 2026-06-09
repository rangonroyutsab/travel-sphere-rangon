package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	"travel-sphere-rangon/models"

	beego "github.com/beego/beego/v2/server/web"
)

const (
	defaultAttractionRadius = "50000"
	defaultAttractionLimit  = "12"
	defaultAttractionFormat = "json"
	defaultAttractionKinds  = "interesting_places,museums,historic,architecture,natural"
)

type OpenTripMapClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

func NewOpenTripMapClient() *OpenTripMapClient {
	baseURL, _ := beego.AppConfig.String("OPENTRIPMAP_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.opentripmap.com/0.1/en/places"
	}
	apiKey, _ := beego.AppConfig.String("OPENTRIPMAP_API_KEY")

	return &OpenTripMapClient{
		BaseURL: strings.TrimRight(baseURL, "/"),
		APIKey:  strings.TrimSpace(apiKey),
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type openTripMapPlaceResponse struct {
	XID   string  `json:"xid"`
	Name  string  `json:"name"`
	Kinds string  `json:"kinds"`
	Dist  float64 `json:"dist"`
	Point struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"point"`
}

func (c *OpenTripMapClient) GetAttractionsByCoordinates(lat, lng float64) ([]models.Attraction, error) {
	if err := ValidateCoordinates(lat, lng); err != nil {
		return nil, err
	}
	if c.APIKey == "" {
		return []models.Attraction{}, errors.New("opentripmap api key is missing")
	}
	endpoint, err := url.Parse(c.BaseURL + "/radius")
	if err != nil {
		return nil, err
	}

	query := endpoint.Query()
	query.Set("radius", defaultAttractionRadius)
	query.Set("lon", fmt.Sprintf("%f", lng))
	query.Set("lat", fmt.Sprintf("%f", lat))
	query.Set("limit", defaultAttractionLimit)
	query.Set("format", defaultAttractionFormat)
	query.Set("apikey", c.APIKey)
	query.Set("kinds", defaultAttractionKinds)
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("opentripmap returned status %d", res.StatusCode)
	}

	var raw []openTripMapPlaceResponse
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		return nil, err
	}

	attractions := make([]models.Attraction, 0, len(raw))
	for _, item := range raw {
		if strings.TrimSpace(item.Name) == "" {
			continue
		}
		attractions = append(attractions, models.Attraction{
			ID:       item.XID,
			Name:     item.Name,
			Kind:     formatKinds(item.Kinds),
			Distance: item.Dist,
			Lat:      item.Point.Lat,
			Lng:      item.Point.Lon,
		})
	}
	return attractions, nil
}

func formatKinds(kinds string) string {
	if strings.TrimSpace(kinds) == "" {
		return ""
	}

	parts := strings.Split(kinds, ",")

	for i, part := range parts {
		parts[i] = strings.ReplaceAll(strings.TrimSpace(part), "_", " ")
	}

	return strings.Join(parts, ", ")
}
