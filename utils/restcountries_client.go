package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
	"travel-sphere-rangon/models"

	beego "github.com/beego/beego/v2/server/web"
)

type RestCountriesClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewRestCountriesClient() *RestCountriesClient {
	baseURL, _ := beego.AppConfig.String("REST_COUNTRIES_BASE_URL")
	if baseURL == "" {
		baseURL = "https://restcountries.com/v3.1"
	}
	return &RestCountriesClient{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type restCountryResponse struct {
	Name struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	} `json:"name"`

	Capital   []string `json:"capital"`
	Region    string   `json:"region"`
	Subregion string   `json:"subregion"`

	Population int `json:"population"`

	Flags struct {
		PNG string `json:"png"`
		SVG string `json:"svg"`
	} `json:"flags"`

	Currencies map[string]struct {
		Name string `json:"name"`
	} `json:"currencies"`

	Languages map[string]string `json:"languages"`

	LatLng []float64 `json:"latlng"`
}

func (c *RestCountriesClient) GetAllCountries() ([]models.Country, error) {
	url := c.BaseURL + "/all?fields=name,capital,region,subregion,population,flags,currencies,languages,latlng"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("rest countries returned status %d", res.StatusCode)
	}

	var raw []restCountryResponse
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		return nil, err
	}

	countries := make([]models.Country, 0, len(raw))

	for _, item := range raw {
		countries = append(countries, transformCountry(item))
	}

	sort.Slice(countries, func(i, j int) bool {
		return countries[i].Name < countries[j].Name
	})

	return countries, nil
}

func transformCountry(item restCountryResponse) models.Country {
	capital := "N/A"
	if len(item.Capital) > 0 && strings.TrimSpace(item.Capital[0]) != "" {
		capital = item.Capital[0]
	}

	flag := item.Flags.PNG
	if flag == "" {
		flag = item.Flags.SVG
	}

	currency := "N/A"
	for code, cur := range item.Currencies {
		if cur.Name != "" {
			currency = fmt.Sprintf("%s (%s)", code, cur.Name)
		} else {
			currency = code
		}
		break
	}

	languages := make([]string, 0, len(item.Languages))
	for _, language := range item.Languages {
		languages = append(languages, language)
	}
	sort.Strings(languages)

	lat := 0.0
	lng := 0.0
	if len(item.LatLng) >= 2 {
		lat = item.LatLng[0]
		lng = item.LatLng[1]
	}

	return models.Country{
		Name:       item.Name.Common,
		Official:   item.Name.Official,
		Slug:       Slugify(item.Name.Common),
		Capital:    capital,
		Region:     item.Region,
		Subregion:  item.Subregion,
		Population: FormatPopulation(item.Population),
		Flag:       flag,
		Currency:   currency,
		Languages:  languages,
		Lat:        lat,
		Lng:        lng,
	}
}
