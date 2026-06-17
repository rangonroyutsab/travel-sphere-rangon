package services

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

const testRestCountriesResponseFields = "names,capitals,region,subregion,population,flag,currencies,languages,coordinates"

func newMockCountryService(t *testing.T, statusCode int, responseBody string) (*CountryService, func()) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertRestCountriesServiceRequest(t, r, 0)

		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(responseBody))
	}))

	service := NewCountryService()
	service.client.BaseURL = server.URL
	service.client.APIKey = "test-api-key"

	return service, server.Close
}

func assertRestCountriesServiceRequest(t *testing.T, r *http.Request, expectedOffset int) {
	t.Helper()

	if r.URL.Path != "/" {
		t.Fatalf("expected root path, got %s", r.URL.Path)
	}

	if r.Header.Get("Authorization") != "Bearer test-api-key" {
		t.Fatalf("expected bearer auth header, got %q", r.Header.Get("Authorization"))
	}

	if r.Header.Get("Accept") != "application/json" {
		t.Fatalf("expected Accept application/json, got %q", r.Header.Get("Accept"))
	}

	query := r.URL.Query()
	if query.Get("response_fields") != testRestCountriesResponseFields {
		t.Fatalf("expected response_fields %q, got %q", testRestCountriesResponseFields, query.Get("response_fields"))
	}

	if query.Get("limit") != "100" {
		t.Fatalf("expected limit 100, got %q", query.Get("limit"))
	}

	if query.Get("offset") != strconv.Itoa(expectedOffset) {
		t.Fatalf("expected offset %d, got %q", expectedOffset, query.Get("offset"))
	}
}

func restCountriesServicePage(objects string, count int, offset int, more bool) string {
	return fmt.Sprintf(`{
		"data": {
			"objects": [%s],
			"meta": {
				"count": %d,
				"offset": %d,
				"more": %t
			}
		}
	}`, objects, count, offset, more)
}

func TestCountryServiceGetAllCountries(t *testing.T) {
	service, cleanup := newMockCountryService(t, http.StatusOK, restCountriesServicePage(`{
		"names": {
			"common": "Japan",
			"official": "Japan"
		},
		"capitals": [{"name": "Tokyo"}],
		"region": "Asia",
		"subregion": "Eastern Asia",
		"population": 125800000,
		"flag": {
			"url_png": "https://example.com/japan.png",
			"url_svg": "https://example.com/japan.svg"
		},
		"currencies": [
			{
				"code": "JPY",
				"name": "Japanese yen"
			}
		],
		"languages": [
			{"name": "Japanese"}
		],
		"coordinates": {
			"lat": 36.0,
			"lng": 138.0
		}
	}`, 1, 0, false))
	defer cleanup()

	countries, err := service.GetAllCountries()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(countries) != 1 {
		t.Fatalf("expected 1 country, got %d", len(countries))
	}

	country := countries[0]

	if country.Name != "Japan" {
		t.Errorf("expected Japan, got %s", country.Name)
	}

	if country.Capital != "Tokyo" {
		t.Errorf("expected Tokyo, got %s", country.Capital)
	}

	if country.Region != "Asia" {
		t.Errorf("expected Asia, got %s", country.Region)
	}

	if country.Population != "125.8M" {
		t.Errorf("expected 125.8M, got %s", country.Population)
	}

	if country.Currency != "JPY (Japanese yen)" {
		t.Errorf("expected JPY currency, got %s", country.Currency)
	}

	if len(country.Languages) != 1 || country.Languages[0] != "Japanese" {
		t.Errorf("expected Japanese language, got %v", country.Languages)
	}

	if country.Slug != "japan" {
		t.Errorf("expected slug japan, got %s", country.Slug)
	}

	if country.Lat != 36.0 {
		t.Errorf("expected latitude 36.0, got %f", country.Lat)
	}

	if country.Lng != 138.0 {
		t.Errorf("expected longitude 138.0, got %f", country.Lng)
	}
}

func TestCountryServiceSearchCountries(t *testing.T) {
	service, cleanup := newMockCountryService(t, http.StatusOK, restCountriesServicePage(`{
		"names": {
			"common": "Japan",
			"official": "Japan"
		},
		"capitals": [{"name": "Tokyo"}],
		"region": "Asia",
		"subregion": "Eastern Asia",
		"population": 125800000,
		"flag": {
			"url_png": "https://example.com/japan.png"
		},
		"currencies": [
			{
				"code": "JPY",
				"name": "Japanese yen"
			}
		],
		"languages": [
			{"name": "Japanese"}
		],
		"coordinates": {
			"lat": 36.0,
			"lng": 138.0
		}
	},
	{
		"names": {
			"common": "France",
			"official": "French Republic"
		},
		"capitals": [{"name": "Paris"}],
		"region": "Europe",
		"subregion": "Western Europe",
		"population": 67000000,
		"flag": {
			"url_png": "https://example.com/france.png"
		},
		"currencies": [
			{
				"code": "EUR",
				"name": "Euro"
			}
		],
		"languages": [
			{"name": "French"}
		],
		"coordinates": {
			"lat": 46.0,
			"lng": 2.0
		}
	}`, 2, 0, false))
	defer cleanup()

	tests := []struct {
		name          string
		search        string
		region        string
		expectedCount int
		expectedName  string
	}{
		{
			name:          "search by country name",
			search:        "japan",
			region:        "",
			expectedCount: 1,
			expectedName:  "Japan",
		},
		{
			name:          "search by capital",
			search:        "paris",
			region:        "",
			expectedCount: 1,
			expectedName:  "France",
		},
		{
			name:          "filter by region",
			search:        "",
			region:        "Europe",
			expectedCount: 1,
			expectedName:  "France",
		},
		{
			name:          "search and region together",
			search:        "tokyo",
			region:        "Asia",
			expectedCount: 1,
			expectedName:  "Japan",
		},
		{
			name:          "no match",
			search:        "canada",
			region:        "",
			expectedCount: 0,
			expectedName:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			countries, err := service.SearchCountries(tt.search, tt.region)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if len(countries) != tt.expectedCount {
				t.Fatalf("expected %d countries, got %d", tt.expectedCount, len(countries))
			}

			if tt.expectedCount > 0 && countries[0].Name != tt.expectedName {
				t.Errorf("expected %s, got %s", tt.expectedName, countries[0].Name)
			}
		})
	}
}

func TestCountryServiceGetCountryBySlug(t *testing.T) {
	service, cleanup := newMockCountryService(t, http.StatusOK, restCountriesServicePage(`{
		"names": {
			"common": "Bangladesh",
			"official": "People's Republic of Bangladesh"
		},
		"capitals": [{"name": "Dhaka"}],
		"region": "Asia",
		"subregion": "Southern Asia",
		"population": 170000000,
		"flag": {
			"url_png": "https://example.com/bangladesh.png"
		},
		"currencies": [
			{
				"code": "BDT",
				"name": "Bangladeshi taka"
			}
		],
		"languages": [
			{"name": "Bengali"}
		],
		"coordinates": {
			"lat": 24.0,
			"lng": 90.0
		}
	}`, 1, 0, false))
	defer cleanup()

	country, err := service.GetCountryBySlug("bangladesh")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if country.Name != "Bangladesh" {
		t.Errorf("expected Bangladesh, got %s", country.Name)
	}
}

func TestCountryServiceGetCountryBySlugReturnsErrorForMissingSlug(t *testing.T) {
	service, cleanup := newMockCountryService(t, http.StatusOK, restCountriesServicePage(`{
		"names": {
			"common": "Bangladesh",
			"official": "People's Republic of Bangladesh"
		},
		"capitals": [{"name": "Dhaka"}],
		"region": "Asia",
		"subregion": "Southern Asia",
		"population": 170000000,
		"flag": {
			"url_png": "https://example.com/bangladesh.png"
		},
		"currencies": [
			{
				"code": "BDT",
				"name": "Bangladeshi taka"
			}
		],
		"languages": [
			{"name": "Bengali"}
		],
		"coordinates": {
			"lat": 24.0,
			"lng": 90.0
		}
	}`, 1, 0, false))
	defer cleanup()

	_, err := service.GetCountryBySlug("missing-country")
	if err == nil {
		t.Fatal("expected error for missing country slug")
	}
}

func TestCountryServiceGetDefaultCountriesLimitsToTwentyFour(t *testing.T) {
	response := ""

	for i := 1; i <= 30; i++ {
		if i > 1 {
			response += `,`
		}

		response += fmt.Sprintf(`{
			"names": {
				"common": "Country %02d",
				"official": "Country %02d Official"
			},
			"capitals": [{"name": "Capital"}],
			"region": "Asia",
			"subregion": "Test",
			"population": 1000,
			"flag": {
				"url_png": "https://example.com/flag.png"
			},
			"currencies": [
				{
					"code": "USD",
					"name": "Dollar"
				}
			],
			"languages": [
				{"name": "English"}
			],
			"coordinates": {
				"lat": 1.0,
				"lng": 2.0
			}
		}`, i, i)
	}

	service, cleanup := newMockCountryService(t, http.StatusOK, restCountriesServicePage(response, 30, 0, false))
	defer cleanup()

	countries, err := service.GetDefaultCountries()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(countries) != 24 {
		t.Fatalf("expected 24 countries, got %d", len(countries))
	}
}

func TestCountryServiceHandlesMissingOptionalFields(t *testing.T) {
	service, cleanup := newMockCountryService(t, http.StatusOK, restCountriesServicePage(`{
		"names": {
			"common": "Testland",
			"official": "Republic of Testland"
		},
		"capitals": [],
		"region": "Test Region",
		"subregion": "",
		"population": 999,
		"flag": {},
		"currencies": [],
		"languages": []
	}`, 1, 0, false))
	defer cleanup()

	countries, err := service.GetAllCountries()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	country := countries[0]

	if country.Capital != "N/A" {
		t.Errorf("expected capital N/A, got %s", country.Capital)
	}

	if country.Currency != "N/A" {
		t.Errorf("expected currency N/A, got %s", country.Currency)
	}

	if len(country.Languages) != 0 {
		t.Errorf("expected no languages, got %v", country.Languages)
	}

	if country.Lat != 0 {
		t.Errorf("expected latitude 0, got %f", country.Lat)
	}

	if country.Lng != 0 {
		t.Errorf("expected longitude 0, got %f", country.Lng)
	}
}

func TestCountryServiceReturnsErrorWhenAPIResponseFails(t *testing.T) {
	service, cleanup := newMockCountryService(t, http.StatusInternalServerError, `{"message":"server error"}`)
	defer cleanup()

	_, err := service.GetAllCountries()
	if err == nil {
		t.Fatal("expected error for failed API response")
	}
}

func TestCountryServiceReturnsErrorForInvalidJSON(t *testing.T) {
	service, cleanup := newMockCountryService(t, http.StatusOK, `{invalid json}`)
	defer cleanup()

	_, err := service.GetAllCountries()
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
