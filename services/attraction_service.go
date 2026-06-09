package services

import (
	"travel-sphere-rangon/models"
	"travel-sphere-rangon/utils"
)

type AttractionService struct {
	client *utils.OpenTripMapClient
}

func NewAttractionService() *AttractionService {
	return &AttractionService{
		client: utils.NewOpenTripMapClient(),
	}
}

func (s *AttractionService) GetAttractions(lat, lng float64) []models.Attraction {
	attractions, err := s.client.GetAttractionsByCoordinates(lat, lng)
	if err != nil {
		return []models.Attraction{}
	}
	return attractions
}
