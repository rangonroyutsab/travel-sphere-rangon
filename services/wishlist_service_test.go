package services

import (
	"testing"

	"travel-sphere-rangon/models"
)

func TestWishlistServiceCreateAndList(t *testing.T) {
	service := NewWishlistService()

	item, err := service.Create("user1", "Japan", "Visit Tokyo", models.StatusPlanned)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if item.ID == "" {
		t.Fatal("expected item ID to be set")
	}

	if item.Username != "user1" {
		t.Errorf("expected username user1, got %s", item.Username)
	}

	if item.CountryName != "Japan" {
		t.Errorf("expected country Japan, got %s", item.CountryName)
	}

	if item.Note != "Visit Tokyo" {
		t.Errorf("expected note Visit Tokyo, got %s", item.Note)
	}

	if item.Status != models.StatusPlanned {
		t.Errorf("expected status Planned, got %s", item.Status)
	}

	items := service.List("user1")

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	if items[0].CountryName != "Japan" {
		t.Errorf("expected Japan, got %s", items[0].CountryName)
	}
}

func TestWishlistServiceDefaultStatusIsPlanned(t *testing.T) {
	service := NewWishlistService()

	item, err := service.Create("user1", "France", "", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if item.Status != models.StatusPlanned {
		t.Errorf("expected default status Planned, got %s", item.Status)
	}
}

func TestWishlistServiceUpdate(t *testing.T) {
	service := NewWishlistService()

	item, err := service.Create("user1", "Brazil", "Old note", models.StatusPlanned)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, err := service.Update("user1", item.ID, "New note", models.StatusVisited)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updated.Note != "New note" {
		t.Errorf("expected note New note, got %s", updated.Note)
	}

	if updated.Status != models.StatusVisited {
		t.Errorf("expected status Visited, got %s", updated.Status)
	}
}

func TestWishlistServiceDelete(t *testing.T) {
	service := NewWishlistService()

	item, err := service.Create("user1", "Argentina", "", models.StatusPlanned)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := service.Delete("user1", item.ID); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	items := service.List("user1")

	if len(items) != 0 {
		t.Fatalf("expected 0 items after delete, got %d", len(items))
	}
}

func TestWishlistServiceKeepsUsersSeparated(t *testing.T) {
	service := NewWishlistService()

	_, err := service.Create("user1", "Japan", "", models.StatusPlanned)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = service.Create("user2", "France", "", models.StatusVisited)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	user1Items := service.List("user1")
	user2Items := service.List("user2")

	if len(user1Items) != 1 {
		t.Fatalf("expected user1 to have 1 item, got %d", len(user1Items))
	}

	if len(user2Items) != 1 {
		t.Fatalf("expected user2 to have 1 item, got %d", len(user2Items))
	}

	if user1Items[0].CountryName != "Japan" {
		t.Errorf("expected user1 item Japan, got %s", user1Items[0].CountryName)
	}

	if user2Items[0].CountryName != "France" {
		t.Errorf("expected user2 item France, got %s", user2Items[0].CountryName)
	}
}

func TestWishlistServiceReturnsEmptyListForUnknownOrEmptyUser(t *testing.T) {
	service := NewWishlistService()

	tests := []struct {
		name     string
		username string
	}{
		{name: "unknown user", username: "missing-user"},
		{name: "empty user", username: ""},
		{name: "spaces only user", username: "   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := service.List(tt.username)

			if len(items) != 0 {
				t.Fatalf("expected empty list, got %d items", len(items))
			}
		})
	}
}

func TestWishlistServiceCreateValidation(t *testing.T) {
	service := NewWishlistService()

	tests := []struct {
		name        string
		username    string
		countryName string
		status      string
	}{
		{
			name:        "empty username",
			username:    "",
			countryName: "Japan",
			status:      models.StatusPlanned,
		},
		{
			name:        "empty country name",
			username:    "user1",
			countryName: "",
			status:      models.StatusPlanned,
		},
		{
			name:        "invalid status",
			username:    "user1",
			countryName: "Japan",
			status:      "Maybe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Create(tt.username, tt.countryName, "", tt.status)

			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestWishlistServiceUpdateValidation(t *testing.T) {
	service := NewWishlistService()

	item, err := service.Create("user1", "Japan", "", models.StatusPlanned)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	tests := []struct {
		name     string
		username string
		id       string
		status   string
	}{
		{
			name:     "empty username",
			username: "",
			id:       item.ID,
			status:   models.StatusVisited,
		},
		{
			name:     "empty id",
			username: "user1",
			id:       "",
			status:   models.StatusVisited,
		},
		{
			name:     "invalid status",
			username: "user1",
			id:       item.ID,
			status:   "Maybe",
		},
		{
			name:     "missing item",
			username: "user1",
			id:       "999",
			status:   models.StatusVisited,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Update(tt.username, tt.id, "Updated note", tt.status)

			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestWishlistServiceDeleteValidation(t *testing.T) {
	service := NewWishlistService()

	tests := []struct {
		name     string
		username string
		id       string
	}{
		{
			name:     "empty username",
			username: "",
			id:       "1",
		},
		{
			name:     "empty id",
			username: "user1",
			id:       "",
		},
		{
			name:     "missing item",
			username: "user1",
			id:       "999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Delete(tt.username, tt.id)

			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}
