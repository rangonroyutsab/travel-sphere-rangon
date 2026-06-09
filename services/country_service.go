package services

import (
	"errors"
	"strings"
	"sync"
	"time"

	"travel-sphere-rangon/models"
	"travel-sphere-rangon/utils"
)

type CountryService struct {
	client     *utils.RestCountriesClient
	cache      []models.Country
	cacheUntil time.Time
	mu         sync.RWMutex
}

func NewCountryService() *CountryService {
	return &CountryService{
		client: utils.NewRestCountriesClient(),
	}
}

func (s *CountryService) GetAllCountries() ([]models.Country, error) {
	s.mu.RLock()

	if time.Now().Before(s.cacheUntil) && len(s.cache) > 0 {
		cached := make([]models.Country, len(s.cache))
		copy(cached, s.cache)

		s.mu.RUnlock()
		return cached, nil
	}

	s.mu.RUnlock()

	countries, err := s.client.GetAllCountries()
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.cache = countries
	s.cacheUntil = time.Now().Add(30 * time.Minute)
	s.mu.Unlock()

	return countries, nil
}

func (s *CountryService) GetDefaultCountries() ([]models.Country, error) {
	countries, err := s.GetAllCountries()
	if err != nil {
		return nil, err
	}

	if len(countries) > 24 {
		return countries[:24], nil
	}

	return countries, nil
}

func (s *CountryService) GetFeaturedCountries() ([]models.Country, error) {
	countries, err := s.GetAllCountries()
	if err != nil {
		return nil, err
	}

	featuredSlugs := []string{
		"united-states",
		"france",
		"japan",
		"australia",
		"brazil",
		"bangladesh",
		"germany",
		"argentina",
	}

	bySlug := make(map[string]models.Country, len(countries))
	for _, country := range countries {
		bySlug[country.Slug] = country
	}

	featured := make([]models.Country, 0, len(featuredSlugs))
	for _, slug := range featuredSlugs {
		country, ok := bySlug[slug]
		if ok {
			featured = append(featured, country)
		}
	}

	return featured, nil
}

func (s *CountryService) SearchCountries(search string, region string) ([]models.Country, error) {
	countries, err := s.GetAllCountries()
	if err != nil {
		return nil, err
	}

	search = strings.ToLower(strings.TrimSpace(search))
	region = strings.ToLower(strings.TrimSpace(region))

	filtered := make([]models.Country, 0)

	for _, country := range countries {
		matchesSearch := true
		matchesRegion := true

		if search != "" {
			name := strings.ToLower(country.Name)
			capital := strings.ToLower(country.Capital)
			official := strings.ToLower(country.Official)

			matchesSearch = strings.Contains(name, search) ||
				strings.Contains(capital, search) ||
				strings.Contains(official, search)
		}

		if region != "" {
			matchesRegion = strings.ToLower(country.Region) == region
		}

		if matchesSearch && matchesRegion {
			filtered = append(filtered, country)
		}
	}

	return filtered, nil
}

func (s *CountryService) GetCountryBySlug(slug string) (*models.Country, error) {
	countries, err := s.GetAllCountries()
	if err != nil {
		return nil, err
	}

	slug = strings.ToLower(strings.TrimSpace(slug))

	for _, country := range countries {
		if country.Slug == slug {
			return &country, nil
		}
	}

	return nil, errors.New("country not found")
}
