package utils

import "testing"

func TestFormatPopulation(t *testing.T) {
	tests := []struct {
		name       string
		population int
		expected   string
	}{
		{
			name:       "small number",
			population: 999,
			expected:   "999",
		},
		{
			name:       "exact thousand",
			population: 1000,
			expected:   "1.0K",
		},
		{
			name:       "thousands",
			population: 1500,
			expected:   "1.5K",
		},
		{
			name:       "exact million",
			population: 1000000,
			expected:   "1.0M",
		},
		{
			name:       "millions",
			population: 2500000,
			expected:   "2.5M",
		},
		{
			name:       "exact billion",
			population: 1000000000,
			expected:   "1.0B",
		},
		{
			name:       "billions",
			population: 1700000000,
			expected:   "1.7B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatPopulation(tt.population)

			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestJoinStrings(t *testing.T) {
	tests := []struct {
		name     string
		values   []string
		expected string
	}{
		{
			name:     "two values",
			values:   []string{"Bangla", "English"},
			expected: "Bangla, English",
		},
		{
			name:     "one value",
			values:   []string{"Japanese"},
			expected: "Japanese",
		},
		{
			name:     "empty slice",
			values:   []string{},
			expected: "N/A",
		},
		{
			name:     "nil slice",
			values:   nil,
			expected: "N/A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := JoinStrings(tt.values)

			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple country",
			input:    "Japan",
			expected: "japan",
		},
		{
			name:     "country with spaces",
			input:    "United States",
			expected: "united-states",
		},
		{
			name:     "country with surrounding spaces",
			input:    "  New Zealand  ",
			expected: "new-zealand",
		},
		{
			name:     "apostrophe removed",
			input:    "Côte d'Ivoire",
			expected: "côte-divoire",
		},
		{
			name:     "ampersand replaced",
			input:    "Bosnia & Herzegovina",
			expected: "bosnia-and-herzegovina",
		},
		{
			name:     "period removed",
			input:    "St. Martin",
			expected: "st-martin",
		},
		{
			name:     "comma removed",
			input:    "Korea, Republic of",
			expected: "korea-republic-of",
		},
		{
			name:     "parentheses removed",
			input:    "Bolivia (Plurinational State of)",
			expected: "bolivia-plurinational-state-of",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Slugify(tt.input)

			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}
