package character

import (
	"testing"
	characterRanking "wowperf/internal/services/warcraftlogs/character"
)

func TestDecodeCharacterName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "Simple name without encoding",
			input:    "Simple",
			expected: "Simple",
			wantErr:  false,
		},
		{
			name:     "Name with single encoding",
			input:    "Lagar%C3%B8",
			expected: "Lagarø",
			wantErr:  false,
		},
		{
			name:     "Name with double encoding",
			input:    "Lagar%25C3%25B8",
			expected: "Lagarø",
			wantErr:  false,
		},
		{
			name:     "Name with special characters",
			input:    "T%C3%A9st%C3%B8r",
			expected: "Téstør",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := characterRanking.DecodeCharacterName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeCharacterName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("DecodeCharacterName() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Test the character ranking query
// go test ./internal/services/warcraftlogs/character/test -v
