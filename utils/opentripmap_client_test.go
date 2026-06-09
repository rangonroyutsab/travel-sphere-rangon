package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newMockOpenTripMapClient(t *testing.T, statusCode int, responseBody string) (*OpenTripMapClient, func()) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/radius" {
			t.Fatalf("expected path /radius, got %s", r.URL.Path)
		}

		if r.URL.Query().Get("radius") == "" {
			t.Fatal("expected radius query parameter")
		}

		if r.URL.Query().Get("lat") == "" {
			t.Fatal("expected lat query parameter")
		}

		if r.URL.Query().Get("lon") == "" {
			t.Fatal("expected lon query parameter")
		}

		if r.URL.Query().Get("limit") == "" {
			t.Fatal("expected limit query parameter")
		}

		if r.URL.Query().Get("format") != "json" {
			t.Fatalf("expected format json, got %s", r.URL.Query().Get("format"))
		}

		if r.URL.Query().Get("apikey") == "" {
			t.Fatal("expected apikey query parameter")
		}

		if r.URL.Query().Get("kinds") == "" {
			t.Fatal("expected kinds query parameter")
		}

		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(responseBody))
	}))

	client := &OpenTripMapClient{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
		HTTPClient: &http.Client{
			Timeout: 2 * time.Second,
		},
	}

	return client, server.Close
}

func TestOpenTripMapClientGetAttractionsByCoordinates(t *testing.T) {
	client, cleanup := newMockOpenTripMapClient(t, http.StatusOK, `[
		{
			"xid": "A1",
			"name": "Lalbagh Fort",
			"kinds": "historic,architecture",
			"dist": 1200.5,
			"point": {
				"lat": 23.7189,
				"lon": 90.3883
			}
		},
		{
			"xid": "A2",
			"name": "Ahsan Manzil",
			"kinds": "museums,interesting_places",
			"dist": 2300.25,
			"point": {
				"lat": 23.7086,
				"lon": 90.4066
			}
		}
	]`)
	defer cleanup()

	attractions, err := client.GetAttractionsByCoordinates(23.8103, 90.4125)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(attractions) != 2 {
		t.Fatalf("expected 2 attractions, got %d", len(attractions))
	}

	first := attractions[0]

	if first.ID != "A1" {
		t.Errorf("expected ID A1, got %s", first.ID)
	}

	if first.Name != "Lalbagh Fort" {
		t.Errorf("expected Lalbagh Fort, got %s", first.Name)
	}

	if first.Kind != "historic, architecture" {
		t.Errorf("expected formatted kind, got %s", first.Kind)
	}

	if first.Distance != 1200.5 {
		t.Errorf("expected distance 1200.5, got %f", first.Distance)
	}

	if first.Lat != 23.7189 {
		t.Errorf("expected latitude 23.7189, got %f", first.Lat)
	}

	if first.Lng != 90.3883 {
		t.Errorf("expected longitude 90.3883, got %f", first.Lng)
	}
}

func TestOpenTripMapClientSkipsAttractionsWithEmptyName(t *testing.T) {
	client, cleanup := newMockOpenTripMapClient(t, http.StatusOK, `[
		{
			"xid": "A1",
			"name": "",
			"kinds": "historic",
			"dist": 100,
			"point": {
				"lat": 23.7,
				"lon": 90.3
			}
		},
		{
			"xid": "A2",
			"name": "National Museum",
			"kinds": "museums",
			"dist": 200,
			"point": {
				"lat": 23.8,
				"lon": 90.4
			}
		}
	]`)
	defer cleanup()

	attractions, err := client.GetAttractionsByCoordinates(23.8103, 90.4125)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(attractions) != 1 {
		t.Fatalf("expected 1 attraction, got %d", len(attractions))
	}

	if attractions[0].Name != "National Museum" {
		t.Errorf("expected National Museum, got %s", attractions[0].Name)
	}
}

func TestOpenTripMapClientFormatsKinds(t *testing.T) {
	tests := []struct {
		name     string
		kinds    string
		expected string
	}{
		{
			name:     "comma separated kinds",
			kinds:    "historic,architecture",
			expected: "historic, architecture",
		},
		{
			name:     "underscore kind",
			kinds:    "interesting_places,museums",
			expected: "interesting places, museums",
		},
		{
			name:     "empty kind",
			kinds:    "",
			expected: "",
		},
		{
			name:     "spaces around kind",
			kinds:    " historic , tourist_facilities ",
			expected: "historic, tourist facilities",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatKinds(tt.kinds)

			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestOpenTripMapClientReturnsErrorWhenAPIKeyMissing(t *testing.T) {
	client, cleanup := newMockOpenTripMapClient(t, http.StatusOK, `[]`)
	defer cleanup()

	client.APIKey = ""

	_, err := client.GetAttractionsByCoordinates(23.8103, 90.4125)
	if err == nil {
		t.Fatal("expected error for missing API key")
	}
}

func TestOpenTripMapClientReturnsErrorForInvalidCoordinates(t *testing.T) {
	client, cleanup := newMockOpenTripMapClient(t, http.StatusOK, `[]`)
	defer cleanup()

	tests := []struct {
		name string
		lat  float64
		lng  float64
	}{
		{name: "latitude too low", lat: -91, lng: 90},
		{name: "latitude too high", lat: 91, lng: 90},
		{name: "longitude too low", lat: 23, lng: -181},
		{name: "longitude too high", lat: 23, lng: 181},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.GetAttractionsByCoordinates(tt.lat, tt.lng)

			if err == nil {
				t.Fatal("expected error for invalid coordinates")
			}
		})
	}
}

func TestOpenTripMapClientReturnsErrorForNonSuccessStatus(t *testing.T) {
	client, cleanup := newMockOpenTripMapClient(t, http.StatusUnauthorized, `{"message":"invalid api key"}`)
	defer cleanup()

	_, err := client.GetAttractionsByCoordinates(23.8103, 90.4125)
	if err == nil {
		t.Fatal("expected error for non-success status")
	}
}

func TestOpenTripMapClientReturnsErrorForInvalidJSON(t *testing.T) {
	client, cleanup := newMockOpenTripMapClient(t, http.StatusOK, `{invalid json}`)
	defer cleanup()

	_, err := client.GetAttractionsByCoordinates(23.8103, 90.4125)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
