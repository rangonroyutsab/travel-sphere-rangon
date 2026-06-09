package services

import (
	"testing"

	"travel-sphere-rangon/models"
)

func TestDashboardServiceGetSummary(t *testing.T) {
	wishlistService := NewWishlistService()
	dashboardService := NewDashboardService(wishlistService)

	_, err := wishlistService.Create("user1", "Japan", "", models.StatusPlanned)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = wishlistService.Create("user1", "France", "", models.StatusVisited)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = wishlistService.Create("user1", "Brazil", "", models.StatusPlanned)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	summary := dashboardService.GetSummary("user1")

	if summary.TotalSaved != 3 {
		t.Errorf("expected total saved 3, got %d", summary.TotalSaved)
	}

	if summary.Planned != 2 {
		t.Errorf("expected planned 2, got %d", summary.Planned)
	}

	if summary.Visited != 1 {
		t.Errorf("expected visited 1, got %d", summary.Visited)
	}
}

func TestDashboardServiceGetSummaryForEmptyWishlist(t *testing.T) {
	wishlistService := NewWishlistService()
	dashboardService := NewDashboardService(wishlistService)

	summary := dashboardService.GetSummary("user1")

	if summary.TotalSaved != 0 {
		t.Errorf("expected total saved 0, got %d", summary.TotalSaved)
	}

	if summary.Planned != 0 {
		t.Errorf("expected planned 0, got %d", summary.Planned)
	}

	if summary.Visited != 0 {
		t.Errorf("expected visited 0, got %d", summary.Visited)
	}
}

func TestDashboardServiceGetSummaryUsesUsername(t *testing.T) {
	wishlistService := NewWishlistService()
	dashboardService := NewDashboardService(wishlistService)

	_, err := wishlistService.Create("user1", "Japan", "", models.StatusPlanned)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = wishlistService.Create("user2", "France", "", models.StatusVisited)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	user1Summary := dashboardService.GetSummary("user1")
	user2Summary := dashboardService.GetSummary("user2")

	if user1Summary.TotalSaved != 1 {
		t.Errorf("expected user1 total saved 1, got %d", user1Summary.TotalSaved)
	}

	if user1Summary.Planned != 1 {
		t.Errorf("expected user1 planned 1, got %d", user1Summary.Planned)
	}

	if user1Summary.Visited != 0 {
		t.Errorf("expected user1 visited 0, got %d", user1Summary.Visited)
	}

	if user2Summary.TotalSaved != 1 {
		t.Errorf("expected user2 total saved 1, got %d", user2Summary.TotalSaved)
	}

	if user2Summary.Planned != 0 {
		t.Errorf("expected user2 planned 0, got %d", user2Summary.Planned)
	}

	if user2Summary.Visited != 1 {
		t.Errorf("expected user2 visited 1, got %d", user2Summary.Visited)
	}
}

func TestDashboardServiceGetSavedDestinations(t *testing.T) {
	wishlistService := NewWishlistService()
	dashboardService := NewDashboardService(wishlistService)

	_, err := wishlistService.Create("user1", "Japan", "Tokyo trip", models.StatusPlanned)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	destinations := dashboardService.GetSavedDestinations("user1")

	if len(destinations) != 1 {
		t.Fatalf("expected 1 saved destination, got %d", len(destinations))
	}

	if destinations[0].CountryName != "Japan" {
		t.Errorf("expected Japan, got %s", destinations[0].CountryName)
	}

	if destinations[0].Note != "Tokyo trip" {
		t.Errorf("expected note Tokyo trip, got %s", destinations[0].Note)
	}
}

func TestDashboardServiceGetSavedDestinationsForEmptyUser(t *testing.T) {
	wishlistService := NewWishlistService()
	dashboardService := NewDashboardService(wishlistService)

	destinations := dashboardService.GetSavedDestinations("")

	if len(destinations) != 0 {
		t.Fatalf("expected 0 saved destinations, got %d", len(destinations))
	}
}
