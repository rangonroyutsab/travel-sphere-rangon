package services

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newMockCountryService(t *testing.T, statusCode int, responseBody string) (*CountryService, func()) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/all" {
			t.Fatalf("expected path /all, got %s", r.URL.Path)
		}

		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(responseBody))
	}))

	service := NewCountryService()
	service.client.BaseURL = server.URL

	return service, server.Close
}

func TestCountryServiceGetAllCountries(t *testing.T) {
	service, cleanup := newMockCountryService(t, http.StatusOK, `[
		{
			"name": {
				"common": "Japan",
				"official": "Japan"
			},
			"capital": ["Tokyo"],
			"region": "Asia",
			"subregion": "Eastern Asia",
			"population": 125800000,
			"flags": {
				"png": "https://example.com/japan.png",
				"svg": "https://example.com/japan.svg"
			},
			"currencies": {
				"JPY": {
					"name": "Japanese yen"
				}
			},
			"languages": {
				"jpn": "Japanese"
			},
			"latlng": [36.0, 138.0]
		}
	]`)
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
	service, cleanup := newMockCountryService(t, http.StatusOK, `[
		{
			"name": {
				"common": "Japan",
				"official": "Japan"
			},
			"capital": ["Tokyo"],
			"region": "Asia",
			"subregion": "Eastern Asia",
			"population": 125800000,
			"flags": {
				"png": "https://example.com/japan.png"
			},
			"currencies": {
				"JPY": {
					"name": "Japanese yen"
				}
			},
			"languages": {
				"jpn": "Japanese"
			},
			"latlng": [36.0, 138.0]
		},
		{
			"name": {
				"common": "France",
				"official": "French Republic"
			},
			"capital": ["Paris"],
			"region": "Europe",
			"subregion": "Western Europe",
			"population": 67000000,
			"flags": {
				"png": "https://example.com/france.png"
			},
			"currencies": {
				"EUR": {
					"name": "Euro"
				}
			},
			"languages": {
				"fra": "French"
			},
			"latlng": [46.0, 2.0]
		}
	]`)
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
	service, cleanup := newMockCountryService(t, http.StatusOK, `[
		{
			"name": {
				"common": "Bangladesh",
				"official": "People's Republic of Bangladesh"
			},
			"capital": ["Dhaka"],
			"region": "Asia",
			"subregion": "Southern Asia",
			"population": 170000000,
			"flags": {
				"png": "https://example.com/bangladesh.png"
			},
			"currencies": {
				"BDT": {
					"name": "Bangladeshi taka"
				}
			},
			"languages": {
				"ben": "Bengali"
			},
			"latlng": [24.0, 90.0]
		}
	]`)
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
	service, cleanup := newMockCountryService(t, http.StatusOK, `[
		{
			"name": {
				"common": "Bangladesh",
				"official": "People's Republic of Bangladesh"
			},
			"capital": ["Dhaka"],
			"region": "Asia",
			"subregion": "Southern Asia",
			"population": 170000000,
			"flags": {
				"png": "https://example.com/bangladesh.png"
			},
			"currencies": {
				"BDT": {
					"name": "Bangladeshi taka"
				}
			},
			"languages": {
				"ben": "Bengali"
			},
			"latlng": [24.0, 90.0]
		}
	]`)
	defer cleanup()

	_, err := service.GetCountryBySlug("missing-country")
	if err == nil {
		t.Fatal("expected error for missing country slug")
	}
}

func TestCountryServiceGetDefaultCountriesLimitsToTwentyFour(t *testing.T) {
	response := `[`

	for i := 1; i <= 30; i++ {
		if i > 1 {
			response += `,`
		}

		response += fmt.Sprintf(`{
			"name": {
				"common": "Country %02d",
				"official": "Country %02d Official"
			},
			"capital": ["Capital"],
			"region": "Asia",
			"subregion": "Test",
			"population": 1000,
			"flags": {
				"png": "https://example.com/flag.png"
			},
			"currencies": {
				"USD": {
					"name": "Dollar"
				}
			},
			"languages": {
				"eng": "English"
			},
			"latlng": [1.0, 2.0]
		}`, i, i)
	}

	response += `]`

	service, cleanup := newMockCountryService(t, http.StatusOK, response)
	defer cleanup()

	countries, err := service.GetDefaultCountries()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(countries) != 24 {
		t.Fatalf("expected 24 countries, got %d", len(countries))
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
