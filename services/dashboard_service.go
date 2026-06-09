package services

import "travel-sphere-rangon/models"

type DashboardSummary struct {
	TotalSaved int `json:"total_saved"`
	Planned    int `json:"planned"`
	Visited    int `json:"visited"`
}

type DashboardService struct {
	wishlistService *WishlistService
}

func NewDashboardService(wishlistService *WishlistService) *DashboardService {
	return &DashboardService{
		wishlistService: wishlistService,
	}
}

func (s *DashboardService) GetSummary(username string) DashboardSummary {
	items := s.wishlistService.List(username)

	summary := DashboardSummary{
		TotalSaved: len(items),
	}

	for _, item := range items {
		switch item.Status {
		case models.StatusPlanned:
			summary.Planned++
		case models.StatusVisited:
			summary.Visited++
		}
	}

	return summary
}

func (s *DashboardService) GetSavedDestinations(username string) []models.WishlistItem {
	return s.wishlistService.List(username)
}
