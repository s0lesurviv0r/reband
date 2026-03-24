package formats

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/s0lesurviv0r/reband/types"
)

type BC125PYTestSuite struct {
	suite.Suite
}

func TestBC125PYTestSuite(t *testing.T) {
	suite.Run(t, new(BC125PYTestSuite))
}

func (s *BC125PYTestSuite) TestDecode() {
	t := s.T()

	tt := []struct {
		name     string
		input    string
		expected types.Channel
	}{
		{
			name:  "simple",
			input: "1,Channel 1,462.5500,nfm,none,2,unlocked,off\n",
			expected: types.Channel{
				Index:      1,
				Name:       "Channel 1",
				Frequency:  types.Frequency(462550000),
				Modulation: types.ModulationFM,
				Bandwidth:  12500,
				Tone: types.Tone{
					Type:  types.ToneTypeNone,
					Value: 0,
				},
				Delay:    2 * time.Second,
				Lockout:  false,
				Priority: false,
			},
		},
		{
			name:  "ctcss",
			input: "1,Channel 1,462.5500,nfm,ctcss_156.7,2,unlocked,off\n",
			expected: types.Channel{
				Index:      1,
				Name:       "Channel 1",
				Frequency:  types.Frequency(462550000),
				Modulation: types.ModulationFM,
				Bandwidth:  12500,
				Tone: types.Tone{
					Type:  types.ToneTypeCTCSS,
					Value: 1567,
				},
				Delay:    2 * time.Second,
				Lockout:  false,
				Priority: false,
			},
		},
		{
			name:  "lockout",
			input: "1,Channel 1,462.5500,nfm,none,2,locked,off\n",
			expected: types.Channel{
				Index:      1,
				Name:       "Channel 1",
				Frequency:  types.Frequency(462550000),
				Modulation: types.ModulationFM,
				Bandwidth:  12500,
				Tone: types.Tone{
					Type:  types.ToneTypeNone,
					Value: 0,
				},
				Delay:    2 * time.Second,
				Lockout:  true,
				Priority: false,
			},
		},
		{
			name:  "priority",
			input: "1,Channel 1,462.5500,nfm,none,2,unlocked,on\n",
			expected: types.Channel{
				Index:      1,
				Name:       "Channel 1",
				Frequency:  types.Frequency(462550000),
				Modulation: types.ModulationFM,
				Bandwidth:  12500,
				Tone: types.Tone{
					Type:  types.ToneTypeNone,
					Value: 0,
				},
				Delay:    2 * time.Second,
				Lockout:  false,
				Priority: true,
			},
		},
		{
			name:  "wideband_fm",
			input: "1,Channel 1,162.5500,fm,none,2,unlocked,off\n",
			expected: types.Channel{
				Index:      1,
				Name:       "Channel 1",
				Frequency:  types.Frequency(162550000),
				Modulation: types.ModulationFM,
				Bandwidth:  25000,
				Tone: types.Tone{
					Type:  types.ToneTypeNone,
					Value: 0,
				},
				Delay:    2 * time.Second,
				Lockout:  false,
				Priority: false,
			},
		},
	}

	f := NewBC125PY()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			buf.WriteString("Index,Name,Frequency (MHz),Modulation,CTCSS,Delay (sec),Lockout,Priority\n")
			buf.WriteString(tc.input)

			channels, err := f.Decode(buf)
			require.NoError(t, err)

			assert.Len(t, channels, 1)
			assert.Equal(t, tc.expected, channels[0])
		})
	}
}

func (s *BC125PYTestSuite) TestEncode() {
	t := s.T()

	tt := []struct {
		name     string
		input    types.Channel
		expected string
	}{
		{
			name: "narrowband",
			input: types.Channel{
				Index:      1,
				Name:       "Channel 1",
				Frequency:  types.Frequency(462550000),
				Modulation: types.ModulationFM,
				Bandwidth:  12500,
				Tone: types.Tone{
					Type:  types.ToneTypeNone,
					Value: 0,
				},
				Delay:    2 * time.Second,
				Lockout:  false,
				Priority: false,
			},
			expected: "1,Channel 1,462.5500,nfm,none,2,unlocked,off\n",
		},
		{
			name: "wideband",
			input: types.Channel{
				Index:      1,
				Name:       "Channel 1",
				Frequency:  types.Frequency(162550000),
				Modulation: types.ModulationFM,
				Bandwidth:  25000,
				Tone: types.Tone{
					Type:  types.ToneTypeNone,
					Value: 0,
				},
				Delay:    2 * time.Second,
				Lockout:  false,
				Priority: false,
			},
			expected: "1,Channel 1,162.5500,fm,none,2,unlocked,off\n",
		},
		{
			name: "ctcss",
			input: types.Channel{
				Index:      1,
				Name:       "Channel 1",
				Frequency:  types.Frequency(462550000),
				Modulation: types.ModulationFM,
				Bandwidth:  12500,
				Tone: types.Tone{
					Type:  types.ToneTypeCTCSS,
					Value: 1567,
				},
				Delay:    2 * time.Second,
				Lockout:  false,
				Priority: false,
			},
			expected: "1,Channel 1,462.5500,nfm,ctcss_156.7,2,unlocked,off\n",
		},
	}

	f := NewBC125PY()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			encoded := new(bytes.Buffer)
			err := f.Encode(encoded, []types.Channel{tc.input})
			require.NoError(t, err)
			expected := "Index,Name,Frequency (MHz),Modulation,CTCSS,Delay (sec),Lockout,Priority\n"
			expected += tc.expected
			assert.Equal(t, expected, encoded.String())
		})
	}
}
