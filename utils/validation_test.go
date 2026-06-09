package utils

import "testing"

func TestValidateCountryName(t *testing.T) {
	tests := []struct {
		name        string
		countryName string
		wantErr     bool
	}{
		{
			name:        "valid country name",
			countryName: "Japan",
			wantErr:     false,
		},
		{
			name:        "valid country with spaces",
			countryName: "United States",
			wantErr:     false,
		},
		{
			name:        "empty country name",
			countryName: "",
			wantErr:     true,
		},
		{
			name:        "spaces only country name",
			countryName: "   ",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCountryName(tt.countryName)

			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestValidateWishlistStatus(t *testing.T) {
	tests := []struct {
		name    string
		status  string
		wantErr bool
	}{
		{
			name:    "planned status",
			status:  "Planned",
			wantErr: false,
		},
		{
			name:    "visited status",
			status:  "Visited",
			wantErr: false,
		},
		{
			name:    "empty status",
			status:  "",
			wantErr: true,
		},
		{
			name:    "lowercase planned",
			status:  "planned",
			wantErr: true,
		},
		{
			name:    "lowercase visited",
			status:  "visited",
			wantErr: true,
		},
		{
			name:    "invalid status",
			status:  "Maybe",
			wantErr: true,
		},
		{
			name:    "spaces only status",
			status:  "   ",
			wantErr: true,
		},
		{
			name:    "status with surrounding spaces",
			status:  " Planned ",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateWishlistStatus(tt.status)

			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestValidateWishlistID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "valid id",
			id:      "1",
			wantErr: false,
		},
		{
			name:    "valid string id",
			id:      "abc123",
			wantErr: false,
		},
		{
			name:    "empty id",
			id:      "",
			wantErr: true,
		},
		{
			name:    "spaces only id",
			id:      "   ",
			wantErr: true,
		},
		{
			name:    "id with surrounding spaces",
			id:      " 1 ",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateWishlistID(tt.id)

			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestValidateCoordinates(t *testing.T) {
	tests := []struct {
		name    string
		lat     float64
		lng     float64
		wantErr bool
	}{
		{
			name:    "valid Dhaka coordinates",
			lat:     23.8103,
			lng:     90.4125,
			wantErr: false,
		},
		{
			name:    "minimum valid coordinates",
			lat:     -90,
			lng:     -180,
			wantErr: false,
		},
		{
			name:    "maximum valid coordinates",
			lat:     90,
			lng:     180,
			wantErr: false,
		},
		{
			name:    "latitude too low",
			lat:     -91,
			lng:     90,
			wantErr: true,
		},
		{
			name:    "latitude too high",
			lat:     91,
			lng:     90,
			wantErr: true,
		},
		{
			name:    "longitude too low",
			lat:     23,
			lng:     -181,
			wantErr: true,
		},
		{
			name:    "longitude too high",
			lat:     23,
			lng:     181,
			wantErr: true,
		},
		{
			name:    "both invalid",
			lat:     100,
			lng:     200,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCoordinates(tt.lat, tt.lng)

			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}
