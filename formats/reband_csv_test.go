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

type RebandCSVTestSuite struct {
	suite.Suite
}

func TestRebandCSVTestSuite(t *testing.T) {
	suite.Run(t, new(RebandCSVTestSuite))
}

const rebandHeader = "Index,Name,AlphaTag,Comment,Frequency,Duplex,Offset,ToneType,ToneValue,Modulation,Power,Delay,Lockout,Priority\n"

func (s *RebandCSVTestSuite) TestDecode() {
	t := s.T()

	tt := []struct {
		name     string
		input    string
		expected types.Channel
	}{
		{
			name:  "simple",
			input: "1,Channel 1,CH1,A comment,462.5500,,0.0000,none,0,nfm,0,2,false,false\n",
			expected: types.Channel{
				Index:      1,
				Name:       "Channel 1",
				AlphaTag:   "CH1",
				Comment:    "A comment",
				Frequency:  types.Frequency(462550000),
				Modulation: types.ModulationNFM,
				Tone:       types.Tone{Type: types.ToneTypeNone},
				Delay:      2 * time.Second,
			},
		},
		{
			name:  "ctcss",
			input: "2,Repeater,RPT,,146.5200,+,0.6000,ctcss,1567,fm,50,0,false,false\n",
			expected: types.Channel{
				Index:      2,
				Name:       "Repeater",
				AlphaTag:   "RPT",
				Frequency:  types.Frequency(146520000),
				Duplex:     types.DuplexPlus,
				Offset:     types.Frequency(600000),
				Tone:       types.Tone{Type: types.ToneTypeCTCSS, Value: 1567},
				Modulation: types.ModulationFM,
				Power:      50,
			},
		},
		{
			name:  "dcs",
			input: "3,DCS Chan,DCS,,155.3400,-,0.6000,dcs,23,fm,5,0,true,true\n",
			expected: types.Channel{
				Index:      3,
				Name:       "DCS Chan",
				AlphaTag:   "DCS",
				Frequency:  types.Frequency(155340000),
				Duplex:     types.DuplexMinus,
				Offset:     types.Frequency(600000),
				Tone:       types.Tone{Type: types.ToneTypeDCS, Value: 23},
				Modulation: types.ModulationFM,
				Power:      5,
				Lockout:    true,
				Priority:   true,
			},
		},
	}

	f := NewRebandCSV()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			buf.WriteString(rebandHeader)
			buf.WriteString(tc.input)

			channels, err := f.Decode(buf)
			require.NoError(t, err)
			assert.Len(t, channels, 1)
			assert.Equal(t, tc.expected, channels[0])
		})
	}
}

func (s *RebandCSVTestSuite) TestEncode() {
	t := s.T()

	tt := []struct {
		name     string
		input    types.Channel
		expected string
	}{
		{
			name: "simple",
			input: types.Channel{
				Index:      1,
				Name:       "Channel 1",
				AlphaTag:   "CH1",
				Comment:    "A comment",
				Frequency:  types.Frequency(462550000),
				Modulation: types.ModulationNFM,
				Delay:      2 * time.Second,
			},
			expected: "1,Channel 1,CH1,A comment,462.5500,,0.0000,none,0,nfm,0,2,false,false\n",
		},
		{
			name: "ctcss_repeater",
			input: types.Channel{
				Index:      2,
				Name:       "Repeater",
				AlphaTag:   "RPT",
				Frequency:  types.Frequency(146520000),
				Duplex:     types.DuplexPlus,
				Offset:     types.Frequency(600000),
				Tone:       types.Tone{Type: types.ToneTypeCTCSS, Value: 1567},
				Modulation: types.ModulationFM,
				Power:      50,
			},
			expected: "2,Repeater,RPT,,146.5200,+,0.6000,ctcss,1567,fm,50,0,false,false\n",
		},
	}

	f := NewRebandCSV()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := f.Encode(buf, []types.Channel{tc.input})
			require.NoError(t, err)
			assert.Equal(t, rebandHeader+tc.expected, buf.String())
		})
	}
}
