package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrequency_String(t *testing.T) {
	tests := []struct {
		f         int64
		precision int
		str       string
	}{
		{f: 1000000, precision: 4, str: "1.0000"},
		{f: 1000000, precision: 2, str: "1.00"},
		{f: 1000000, precision: 0, str: "1"},
		{f: 1000000000, precision: 4, str: "1000.0000"},
		{f: 1000000000, precision: 2, str: "1000.00"},
		{f: 1000000000, precision: 0, str: "1000"},
	}

	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			f := Frequency(test.f)
			assert.Equal(t, test.str, f.StringWithPrecision(test.precision))
		})
	}
}

func TestTone_CTCSS(t *testing.T) {
	tests := []struct {
		tone Tone
		str  string
	}{
		{tone: Tone{Type: ToneTypeCTCSS, Value: 693}, str: "69.3"},
		{tone: Tone{Type: ToneTypeCTCSS, Value: 1031}, str: "103.1"},
	}

	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			assert.Equal(t, test.str, test.tone.CTCSS())
		})
	}
}
