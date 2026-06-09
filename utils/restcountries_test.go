package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newMockRestCountriesClient(t *testing.T, statusCode int, responseBody string) (*RestCountriesClient, func()) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/all" {
			t.Fatalf("expected path /all, got %s", r.URL.Path)
		}

		if r.URL.Query().Get("fields") == "" {
			t.Fatal("expected fields query parameter")
		}

		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(responseBody))
	}))

	client := &RestCountriesClient{
		BaseURL: server.URL,
		HTTPClient: &http.Client{
			Timeout: 2 * time.Second,
		},
	}

	return client, server.Close
}

func TestRestCountriesClientGetAllCountries(t *testing.T) {
	client, cleanup := newMockRestCountriesClient(t, http.StatusOK, `[
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

	countries, err := client.GetAllCountries()
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

	if country.Official != "Japan" {
		t.Errorf("expected official name Japan, got %s", country.Official)
	}

	if country.Slug != "japan" {
		t.Errorf("expected slug japan, got %s", country.Slug)
	}

	if country.Capital != "Tokyo" {
		t.Errorf("expected capital Tokyo, got %s", country.Capital)
	}

	if country.Region != "Asia" {
		t.Errorf("expected region Asia, got %s", country.Region)
	}

	if country.Subregion != "Eastern Asia" {
		t.Errorf("expected subregion Eastern Asia, got %s", country.Subregion)
	}

	if country.Population != "125.8M" {
		t.Errorf("expected population 125.8M, got %s", country.Population)
	}

	if country.Flag != "https://example.com/japan.png" {
		t.Errorf("expected PNG flag, got %s", country.Flag)
	}

	if country.Currency != "JPY (Japanese yen)" {
		t.Errorf("expected JPY currency, got %s", country.Currency)
	}

	if len(country.Languages) != 1 || country.Languages[0] != "Japanese" {
		t.Errorf("expected Japanese language, got %v", country.Languages)
	}

	if country.Lat != 36.0 {
		t.Errorf("expected latitude 36.0, got %f", country.Lat)
	}

	if country.Lng != 138.0 {
		t.Errorf("expected longitude 138.0, got %f", country.Lng)
	}
}

func TestRestCountriesClientSortsCountriesByName(t *testing.T) {
	client, cleanup := newMockRestCountriesClient(t, http.StatusOK, `[
		{
			"name": {
				"common": "Zimbabwe",
				"official": "Republic of Zimbabwe"
			},
			"capital": ["Harare"],
			"region": "Africa",
			"subregion": "Eastern Africa",
			"population": 15000000,
			"flags": {
				"png": "https://example.com/zimbabwe.png"
			},
			"currencies": {
				"USD": {
					"name": "United States dollar"
				}
			},
			"languages": {
				"eng": "English"
			},
			"latlng": [-20.0, 30.0]
		},
		{
			"name": {
				"common": "Australia",
				"official": "Commonwealth of Australia"
			},
			"capital": ["Canberra"],
			"region": "Oceania",
			"subregion": "Australia and New Zealand",
			"population": 26000000,
			"flags": {
				"png": "https://example.com/australia.png"
			},
			"currencies": {
				"AUD": {
					"name": "Australian dollar"
				}
			},
			"languages": {
				"eng": "English"
			},
			"latlng": [-27.0, 133.0]
		}
	]`)
	defer cleanup()

	countries, err := client.GetAllCountries()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(countries) != 2 {
		t.Fatalf("expected 2 countries, got %d", len(countries))
	}

	if countries[0].Name != "Australia" {
		t.Errorf("expected Australia first, got %s", countries[0].Name)
	}

	if countries[1].Name != "Zimbabwe" {
		t.Errorf("expected Zimbabwe second, got %s", countries[1].Name)
	}
}

func TestRestCountriesClientUsesSVGWhenPNGIsMissing(t *testing.T) {
	client, cleanup := newMockRestCountriesClient(t, http.StatusOK, `[
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
				"svg": "https://example.com/france.svg"
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

	countries, err := client.GetAllCountries()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if countries[0].Flag != "https://example.com/france.svg" {
		t.Errorf("expected SVG flag fallback, got %s", countries[0].Flag)
	}
}

func TestRestCountriesClientHandlesMissingOptionalFields(t *testing.T) {
	client, cleanup := newMockRestCountriesClient(t, http.StatusOK, `[
		{
			"name": {
				"common": "Testland",
				"official": "Republic of Testland"
			},
			"capital": [],
			"region": "Test Region",
			"subregion": "",
			"population": 999,
			"flags": {},
			"currencies": {},
			"languages": {},
			"latlng": []
		}
	]`)
	defer cleanup()

	countries, err := client.GetAllCountries()
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

	if country.Population != "999" {
		t.Errorf("expected population 999, got %s", country.Population)
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

func TestRestCountriesClientReturnsErrorForNonSuccessStatus(t *testing.T) {
	client, cleanup := newMockRestCountriesClient(t, http.StatusInternalServerError, `{"message":"server error"}`)
	defer cleanup()

	_, err := client.GetAllCountries()
	if err == nil {
		t.Fatal("expected error for non-success status")
	}
}

func TestRestCountriesClientReturnsErrorForInvalidJSON(t *testing.T) {
	client, cleanup := newMockRestCountriesClient(t, http.StatusOK, `{invalid json}`)
	defer cleanup()

	_, err := client.GetAllCountries()
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
