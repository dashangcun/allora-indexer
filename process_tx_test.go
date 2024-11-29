package main

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashTx(t *testing.T) {
	tests := []struct {
		name        string
		txRaw       string
		upperCase   bool
		expected    string
		expectError bool
		expectedErr error
	}{
		{
			name:        "Valid Base64, uppercase",
			txRaw:       "dGVzdFR4", // "testTx" in Base64
			upperCase:   true,
			expected:    strings.ToUpper("8C519B7C50548BE414E69281C439E6FD4F6E1DB94D333886BBB523851383251C"),
			expectError: false,
			expectedErr: nil,
		},
		{
			name:        "Valid Base64, lowercase",
			txRaw:       "dGVzdFR4", // "testTx" in Base64
			upperCase:   false,
			expected:    "8c519b7c50548be414e69281c439e6fd4f6e1db94d333886bbb523851383251c",
			expectError: false,
			expectedErr: nil,
		},
		{
			name:        "Invalid Base64",
			txRaw:       "!!!invalid!!!",
			upperCase:   true,
			expected:    "",
			expectError: true,
			expectedErr: base64.CorruptInputError(0),
		},
		{
			name:        "Empty string",
			txRaw:       "",
			upperCase:   false,
			expected:    "",
			expectError: false,
			expectedErr: nil,
		},
		{
			name:        "Actual tx, uppercase",
			txRaw:       "CtwCCtkCCigvZW1pc3Npb25zLnY1Lkluc2VydFdvcmtlclBheWxvYWRSZXF1ZXN0EqwCCithbGxvMXRkZHFlenNyeHMzbDB4ajdyemxrcnQzNWUzeXU3enBxenZqYWxqEvwBCithbGxvMXRkZHFlenNyeHMzbDB4ajdyemxrcnQzNWUzeXU3enBxenZqYWxqEgQIzNVtGBAiPwo9CBAQzNVtGithbGxvMXRkZHFlenNyeHMzbDB4ajdyemxrcnQzNWUzeXU3enBxenZqYWxqIgg5NzE4NC42NipACJ9sncbNjSjGlVPAjsL7pl4+qR6A3LLwMhkfqXnII08IqJXlY+90b3U8st4851ai3HuizlQDOpCgFPn50/purTJCMDNiY2FlMDRjZjBjY2Y4YjA0YjczZGU0MWM0M2EwMWVmZjllMzg1NjI4ZDYyMzMzOWJiMDRhMjM1OTBmM2YwOWQ3Em0KUgpGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQO8rgTPDM+LBLc95BxDoB7/njhWKNYjM5uwSiNZDz8J1xIECgIIARjF0AESFwoRCgV1YWxsbxIIMTY5NjA5MDAQwtEJGkBcner4gHSznxX7DQ6o2P+cWhADmfteRX6v9DcSKRGo3m5wEU6trB7jeL/ttFoyFW58jqlxCyABOh6hsCkf08Q6",
			upperCase:   true,
			expected:    strings.ToUpper("82F2A53227CF363DFFF731BE4A6644E42107886006E2AF192D416860A4CAF960"),
			expectError: false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := hashTx(tt.txRaw, tt.upperCase)
			assert.Equal(t, tt.expectedErr, err)
			if tt.expected == "" {
				assert.Empty(t, result, "Expected empty string")
			} else {
				assert.Equal(t, tt.expected, result, "Hashes do not match")
			}
		})
	}
}
