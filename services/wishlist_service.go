package services

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"travel-sphere-rangon/models"
	"travel-sphere-rangon/utils"
)

type WishlistService struct {
	mu     sync.RWMutex
	items  map[string]map[string]models.WishlistItem
	nextID int
}

func NewWishlistService() *WishlistService {
	return &WishlistService{
		items:  make(map[string]map[string]models.WishlistItem),
		nextID: 1,
	}
}

func (s *WishlistService) List(username string) []models.WishlistItem {
	username = strings.TrimSpace(username)
	if username == "" {
		return []models.WishlistItem{}
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	userItems := s.items[username]
	items := make([]models.WishlistItem, 0, len(userItems))

	for _, item := range userItems {
		items = append(items, item)
	}

	// sorts wishlist items so the most recently added item appears first
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})
	return items
}

func (s *WishlistService) Create(username string, countryName string, note string, status string) (*models.WishlistItem, error) {
	username = strings.TrimSpace(username)
	countryName = strings.TrimSpace(countryName)
	note = strings.TrimSpace(note)
	status = strings.TrimSpace(status)

	if username == "" {
		return nil, errors.New("username is required")
	}

	if err := utils.ValidateCountryName(countryName); err != nil {
		return nil, err
	}

	if status == "" {
		status = models.StatusPlanned
	}

	if err := utils.ValidateWishlistStatus(status); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[username]; !ok {
		s.items[username] = make(map[string]models.WishlistItem)
	}

	id := strconv.Itoa(s.nextID)
	s.nextID++

	item := models.WishlistItem{
		ID:          id,
		Username:    username,
		CountryName: countryName,
		Note:        note,
		Status:      status,
		CreatedAt:   time.Now(),
	}

	s.items[username][id] = item
	return &item, nil
}

func (s *WishlistService) Update(username string, id string, note string, status string) (*models.WishlistItem, error) {
	username = strings.TrimSpace(username)
	id = strings.TrimSpace(id)
	note = strings.TrimSpace(note)
	status = strings.TrimSpace(status)

	if username == "" {
		return nil, errors.New("username is required")
	}

	if err := utils.ValidateWishlistID(id); err != nil {
		return nil, err
	}

	if err := utils.ValidateWishlistStatus(status); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[username][id]
	if !ok {
		return nil, errors.New("wishlist item not found")
	}

	item.Note = note
	item.Status = status
	s.items[username][id] = item

	return &item, nil
}

func (s *WishlistService) Delete(username string, id string) error {
	username = strings.TrimSpace(username)
	id = strings.TrimSpace(id)

	if username == "" {
		return errors.New("username is required")
	}

	if err := utils.ValidateWishlistID(id); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[username][id]; !ok {
		return errors.New("wishlist item not found")
	}

	delete(s.items[username], id)

	return nil
}
