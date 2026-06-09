package services

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func newMockAttractionService(t *testing.T, statusCode int, responseBody string) (*AttractionService, func()) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/radius" {
			t.Fatalf("expected path /radius, got %s", r.URL.Path)
		}

		if r.URL.Query().Get("lat") == "" {
			t.Fatal("expected lat query parameter")
		}

		if r.URL.Query().Get("lon") == "" {
			t.Fatal("expected lon query parameter")
		}

		if r.URL.Query().Get("apikey") == "" {
			t.Fatal("expected apikey query parameter")
		}

		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(responseBody))
	}))

	service := NewAttractionService()
	service.client.BaseURL = server.URL
	service.client.APIKey = "test-api-key"

	return service, server.Close
}

func TestAttractionServiceGetAttractions(t *testing.T) {
	service, cleanup := newMockAttractionService(t, http.StatusOK, `[
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

	attractions := service.GetAttractions(23.8103, 90.4125)

	if len(attractions) != 2 {
		t.Fatalf("expected 2 attractions, got %d", len(attractions))
	}

	if attractions[0].ID != "A1" {
		t.Errorf("expected ID A1, got %s", attractions[0].ID)
	}

	if attractions[0].Name != "Lalbagh Fort" {
		t.Errorf("expected Lalbagh Fort, got %s", attractions[0].Name)
	}

	if attractions[0].Kind != "historic, architecture" {
		t.Errorf("expected formatted kind, got %s", attractions[0].Kind)
	}

	if attractions[0].Distance != 1200.5 {
		t.Errorf("expected distance 1200.5, got %f", attractions[0].Distance)
	}

	if attractions[0].Lat != 23.7189 {
		t.Errorf("expected latitude 23.7189, got %f", attractions[0].Lat)
	}

	if attractions[0].Lng != 90.3883 {
		t.Errorf("expected longitude 90.3883, got %f", attractions[0].Lng)
	}
}

func TestAttractionServiceSkipsAttractionsWithEmptyName(t *testing.T) {
	service, cleanup := newMockAttractionService(t, http.StatusOK, `[
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

	attractions := service.GetAttractions(23.8103, 90.4125)

	if len(attractions) != 1 {
		t.Fatalf("expected 1 attraction, got %d", len(attractions))
	}

	if attractions[0].Name != "National Museum" {
		t.Errorf("expected National Museum, got %s", attractions[0].Name)
	}
}

func TestAttractionServiceReturnsEmptySliceWhenAPIResponseFails(t *testing.T) {
	service, cleanup := newMockAttractionService(t, http.StatusInternalServerError, `{"message":"server error"}`)
	defer cleanup()

	attractions := service.GetAttractions(23.8103, 90.4125)

	if len(attractions) != 0 {
		t.Fatalf("expected empty attractions, got %d", len(attractions))
	}
}

func TestAttractionServiceReturnsEmptySliceForInvalidJSON(t *testing.T) {
	service, cleanup := newMockAttractionService(t, http.StatusOK, `{invalid json}`)
	defer cleanup()

	attractions := service.GetAttractions(23.8103, 90.4125)

	if len(attractions) != 0 {
		t.Fatalf("expected empty attractions, got %d", len(attractions))
	}
}

func TestAttractionServiceReturnsEmptySliceForInvalidCoordinates(t *testing.T) {
	service, cleanup := newMockAttractionService(t, http.StatusOK, `[]`)
	defer cleanup()

	tests := []struct {
		name string
		lat  float64
		lng  float64
	}{
		{name: "invalid latitude low", lat: -91, lng: 90},
		{name: "invalid latitude high", lat: 91, lng: 90},
		{name: "invalid longitude low", lat: 23, lng: -181},
		{name: "invalid longitude high", lat: 23, lng: 181},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attractions := service.GetAttractions(tt.lat, tt.lng)

			if len(attractions) != 0 {
				t.Fatalf("expected empty attractions, got %d", len(attractions))
			}
		})
	}
}

func TestAttractionServiceReturnsEmptySliceWhenAPIKeyMissing(t *testing.T) {
	service, cleanup := newMockAttractionService(t, http.StatusOK, `[]`)
	defer cleanup()

	service.client.APIKey = ""

	attractions := service.GetAttractions(23.8103, 90.4125)

	if len(attractions) != 0 {
		t.Fatalf("expected empty attractions, got %d", len(attractions))
	}
}
