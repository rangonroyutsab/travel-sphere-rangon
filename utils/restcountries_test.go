package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

const testRestCountriesResponseFields = "names,capitals,region,subregion,population,flag,currencies,languages,coordinates"

func newMockRestCountriesClient(t *testing.T, statusCode int, responseBody string) (*RestCountriesClient, func()) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertRestCountriesRequest(t, r, 0)

		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(responseBody))
	}))

	client := &RestCountriesClient{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
		HTTPClient: &http.Client{
			Timeout: 2 * time.Second,
		},
	}

	return client, server.Close
}

func assertRestCountriesRequest(t *testing.T, r *http.Request, expectedOffset int) {
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

func restCountriesPage(objects string, count int, offset int, more bool) string {
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

func TestRestCountriesClientGetAllCountries(t *testing.T) {
	client, cleanup := newMockRestCountriesClient(t, http.StatusOK, restCountriesPage(`{
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
	client, cleanup := newMockRestCountriesClient(t, http.StatusOK, restCountriesPage(`{
		"names": {
			"common": "Zimbabwe",
			"official": "Republic of Zimbabwe"
		},
		"capitals": [{"name": "Harare"}],
		"region": "Africa",
		"subregion": "Eastern Africa",
		"population": 15000000,
		"flag": {
			"url_png": "https://example.com/zimbabwe.png"
		},
		"currencies": [
			{
				"code": "USD",
				"name": "United States dollar"
			}
		],
		"languages": [
			{"name": "English"}
		],
		"coordinates": {
			"lat": -20.0,
			"lng": 30.0
		}
	},
	{
		"names": {
			"common": "Australia",
			"official": "Commonwealth of Australia"
		},
		"capitals": [{"name": "Canberra"}],
		"region": "Oceania",
		"subregion": "Australia and New Zealand",
		"population": 26000000,
		"flag": {
			"url_png": "https://example.com/australia.png"
		},
		"currencies": [
			{
				"code": "AUD",
				"name": "Australian dollar"
			}
		],
		"languages": [
			{"name": "English"}
		],
		"coordinates": {
			"lat": -27.0,
			"lng": 133.0
		}
	}`, 2, 0, false))
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
	client, cleanup := newMockRestCountriesClient(t, http.StatusOK, restCountriesPage(`{
		"names": {
			"common": "France",
			"official": "French Republic"
		},
		"capitals": [{"name": "Paris"}],
		"region": "Europe",
		"subregion": "Western Europe",
		"population": 67000000,
		"flag": {
			"url_svg": "https://example.com/france.svg"
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
	}`, 1, 0, false))
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
	client, cleanup := newMockRestCountriesClient(t, http.StatusOK, restCountriesPage(`{
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

func TestRestCountriesClientReturnsErrorWhenAPIKeyMissing(t *testing.T) {
	previousAPIKey, _ := beego.AppConfig.String("REST_COUNTRIES_API_KEY")
	if err := beego.AppConfig.Set("REST_COUNTRIES_API_KEY", ""); err != nil {
		t.Fatalf("failed to clear API key config: %v", err)
	}
	defer func() {
		_ = beego.AppConfig.Set("REST_COUNTRIES_API_KEY", previousAPIKey)
	}()

	client := &RestCountriesClient{
		BaseURL: "https://example.test",
		APIKey:  "",
		HTTPClient: &http.Client{
			Timeout: 2 * time.Second,
		},
	}

	_, err := client.GetAllCountries()
	if err == nil {
		t.Fatal("expected error when API key is missing")
	}

	if err.Error() != "REST_COUNTRIES_API_KEY is required" {
		t.Fatalf("expected missing API key error, got %v", err)
	}
}

func TestRestCountriesClientRequestsMultiplePagesWithNextOffset(t *testing.T) {
	requestedOffsets := make([]int, 0, 2)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			t.Fatalf("expected numeric offset, got %q", r.URL.Query().Get("offset"))
		}

		requestedOffsets = append(requestedOffsets, offset)
		assertRestCountriesRequest(t, r, offset)

		switch offset {
		case 0:
			_, _ = w.Write([]byte(restCountriesPage(`{
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
			}`, 1, 0, true)))
		case 1:
			_, _ = w.Write([]byte(restCountriesPage(`{
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
			}`, 1, 1, false)))
		default:
			t.Fatalf("unexpected offset %d", offset)
		}
	}))
	defer server.Close()

	client := &RestCountriesClient{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
		HTTPClient: &http.Client{
			Timeout: 2 * time.Second,
		},
	}

	countries, err := client.GetAllCountries()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(countries) != 2 {
		t.Fatalf("expected 2 countries, got %d", len(countries))
	}

	if !reflect.DeepEqual(requestedOffsets, []int{0, 1}) {
		t.Fatalf("expected offsets 0,1, got %v", requestedOffsets)
	}
}

func TestRestCountriesClientReturnsErrorWhenPaginatingWithZeroCount(t *testing.T) {
	client, cleanup := newMockRestCountriesClient(t, http.StatusOK, restCountriesPage(``, 0, 0, true))
	defer cleanup()

	_, err := client.GetAllCountries()
	if err == nil {
		t.Fatal("expected error for zero-count pagination")
	}

	if err.Error() != "rest countries response cannot paginate with zero count" {
		t.Fatalf("expected zero-count pagination error, got %v", err)
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
