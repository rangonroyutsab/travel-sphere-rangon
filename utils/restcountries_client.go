package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
	"travel-sphere-rangon/models"

	beego "github.com/beego/beego/v2/server/web"
)

type RestCountriesClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

const defaultRestCountriesBaseURL = "https://api.restcountries.com/countries/v5"
const restCountriesPageLimit = 100

func NewRestCountriesClient() *RestCountriesClient {
	baseURL, _ := beego.AppConfig.String("REST_COUNTRIES_BASE_URL")
	if baseURL == "" {
		baseURL = defaultRestCountriesBaseURL
	}

	apiKey, _ := beego.AppConfig.String("REST_COUNTRIES_API_KEY")

	return &RestCountriesClient{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		APIKey:     strings.TrimSpace(apiKey),
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type restCountriesAPIResponse struct {
	Data struct {
		Objects []restCountryResponse `json:"objects"`
		Meta    struct {
			Count  int  `json:"count"`
			Offset int  `json:"offset"`
			More   bool `json:"more"`
		} `json:"meta"`
	} `json:"data"`
}

type restCountryResponse struct {
	Names struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	} `json:"names"`

	Capitals []struct {
		Name string `json:"name"`
	} `json:"capitals"`
	Region    string `json:"region"`
	Subregion string `json:"subregion"`

	Population int `json:"population"`

	Flag struct {
		PNG string `json:"url_png"`
		SVG string `json:"url_svg"`
	} `json:"flag"`

	Currencies []struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"currencies"`

	Languages []struct {
		Name string `json:"name"`
	} `json:"languages"`

	Coordinates *struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"coordinates"`
}

func (c *RestCountriesClient) GetAllCountries() ([]models.Country, error) {
	if strings.TrimSpace(c.APIKey) == "" {
		apiKey, _ := beego.AppConfig.String("REST_COUNTRIES_API_KEY")
		c.APIKey = strings.TrimSpace(apiKey)
		if c.APIKey == "" {
			return nil, errors.New("REST_COUNTRIES_API_KEY is required")
		}
	}

	var raw []restCountryResponse
	offset := 0

	for {
		requestURL, err := c.countriesURL(offset)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodGet, requestURL, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(c.APIKey))
		req.Header.Set("Accept", "application/json")

		res, err := c.HTTPClient.Do(req)
		if err != nil {
			return nil, err
		}

		if res.StatusCode < 200 || res.StatusCode >= 300 {
			res.Body.Close()
			return nil, fmt.Errorf("rest countries returned status %d", res.StatusCode)
		}

		var page restCountriesAPIResponse
		if err := json.NewDecoder(res.Body).Decode(&page); err != nil {
			res.Body.Close()
			return nil, err
		}
		res.Body.Close()

		raw = append(raw, page.Data.Objects...)

		if !page.Data.Meta.More {
			break
		}

		if page.Data.Meta.Count == 0 {
			return nil, errors.New("rest countries response cannot paginate with zero count")
		}

		offset = page.Data.Meta.Offset + page.Data.Meta.Count
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

func (c *RestCountriesClient) countriesURL(offset int) (string, error) {
	parsedURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid REST Countries base URL: %w", err)
	}

	query := parsedURL.Query()
	query.Set("response_fields", "names,capitals,region,subregion,population,flag,currencies,languages,coordinates")
	query.Set("limit", strconv.Itoa(restCountriesPageLimit))
	query.Set("offset", strconv.Itoa(offset))
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}

func transformCountry(item restCountryResponse) models.Country {
	capital := "N/A"
	for _, itemCapital := range item.Capitals {
		if strings.TrimSpace(itemCapital.Name) != "" {
			capital = itemCapital.Name
			break
		}
	}

	flag := item.Flag.PNG
	if flag == "" {
		flag = item.Flag.SVG
	}

	currency := "N/A"
	for _, cur := range item.Currencies {
		code := strings.TrimSpace(cur.Code)
		name := strings.TrimSpace(cur.Name)

		switch {
		case code != "" && name != "":
			currency = fmt.Sprintf("%s (%s)", code, name)
		case code != "":
			currency = code
		case name != "":
			currency = name
		default:
			continue
		}

		break
	}

	languages := make([]string, 0, len(item.Languages))
	for _, language := range item.Languages {
		name := strings.TrimSpace(language.Name)
		if name != "" {
			languages = append(languages, name)
		}
	}
	sort.Strings(languages)

	lat := 0.0
	lng := 0.0
	if item.Coordinates != nil {
		lat = item.Coordinates.Lat
		lng = item.Coordinates.Lng
	}

	return models.Country{
		Name:       item.Names.Common,
		Official:   item.Names.Official,
		Slug:       Slugify(item.Names.Common),
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
