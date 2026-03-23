package formats

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/s0lesurviv0r/channel-conv/types"
)

type ChirpTestSuite struct {
	suite.Suite
}

func TestChirpTestSuite(t *testing.T) {
	suite.Run(t, new(ChirpTestSuite))
}

func (s *ChirpTestSuite) TestDecode() {
	t := s.T()

	tt := []struct {
		name     string
		input    string
		expected types.Channel
	}{
		{
			name:  "simple",
			input: "27,HAILVHF,146.520000,,0.000000,,88.5,88.5,023,NN,023,Tone->Tone,FM,5.00,,50W,,,,,",
			expected: types.Channel{
				Index:      27,
				Name:       "HAILVHF",
				Frequency:  types.Frequency(146520000),
				Modulation: types.ModulationFM,
				Tone:       types.Tone{Type: types.ToneTypeNone},
				Power:      50,
			},
		},
		{
			name:  "repeater",
			input: "35,K6ITR,145.220000,-,0.600000,Tone,103.5,88.5,023,NN,023,Tone->Tone,FM,5.00,,50W,,,,,",
			expected: types.Channel{
				Index:      35,
				Name:       "K6ITR",
				Frequency:  types.Frequency(145220000),
				Modulation: types.ModulationFM,
				Tone:       types.Tone{Type: types.ToneTypeCTCSS, Value: 1035},
				Duplex:     types.DuplexMinus,
				Offset:     types.Frequency(600000),
				Power:      50,
			},
		},
	}

	f := NewChirp()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			buf.WriteString("Location,Name,Frequency,Duplex,Offset,Tone,rToneFreq,cToneFreq,DtcsCode,DtcsPolarity,RxDtcsCode,CrossMode,Mode,TStep,Skip,Power,Comment,URCALL,RPT1CALL,RPT2CALL,DVCODE\n")
			buf.WriteString(tc.input)

			channels, err := f.Decode(buf)
			require.NoError(t, err)

			assert.Len(t, channels, 1)
			assert.Equal(t, tc.expected, channels[0])
		})
	}
}

func (s *ChirpTestSuite) TestEncode() {
	t := s.T()

	tt := []struct {
		name     string
		input    types.Channel
		expected string
	}{
		{
			name: "simple",
			input: types.Channel{
				Index:      27,
				Name:       "HAILVHF",
				Frequency:  types.Frequency(146520000),
				Modulation: types.ModulationFM,
				Tone:       types.Tone{Type: types.ToneTypeNone},
			},
			expected: "27,HAILVHF,146.520000,,0.000000,,88.5,88.5,023,NN,023,Tone->Tone,FM,5.00,,50W,,,,,",
		},
	}

	f := NewChirp()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			encoded := new(bytes.Buffer)
			err := f.Encode(encoded, []types.Channel{tc.input})
			require.NoError(t, err)
			expected := "Location,Name,Frequency,Duplex,Offset,Tone,rToneFreq,cToneFreq,DtcsCode,DtcsPolarity,RxDtcsCode,CrossMode,Mode,TStep,Skip,Power,Comment,URCALL,RPT1CALL,RPT2CALL,DVCODE\n"
			expected += tc.expected
			assert.Equal(t, expected, encoded.String())
		})
	}
}
