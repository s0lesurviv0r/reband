package formats

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/s0lesurviv0r/reband/types"
)

func openSample(t *testing.T, name string) *os.File {
	t.Helper()
	f, err := os.Open("../samples/" + name)
	require.NoError(t, err)
	t.Cleanup(func() { f.Close() })
	return f
}

func TestSampleBC125AT(t *testing.T) {
	f := NewBC125PY()
	channels, err := f.Decode(openSample(t, "bc125at.csv"))
	require.NoError(t, err)
	require.Len(t, channels, 13)

	assert.Equal(t, types.Channel{
		Index:      1,
		Name:       "NOAA WX1",
		Frequency:  162400000,
		Modulation: types.ModulationFM,
		Bandwidth:  25000,
		Tone:       types.Tone{Type: types.ToneTypeNone},
		Delay:      2 * time.Second,
	}, channels[0])

	assert.Equal(t, types.Channel{
		Index:      5,
		Name:       "2m Repeater",
		Frequency:  146940000,
		Modulation: types.ModulationFM,
		Bandwidth:  25000,
		Tone:       types.Tone{Type: types.ToneTypeCTCSS, Value: 1000},
		Delay:      2 * time.Second,
	}, channels[4])

	assert.Equal(t, types.Channel{
		Index:      13,
		Name:       "Locked Out",
		Frequency:  155340000,
		Modulation: types.ModulationFM,
		Bandwidth:  12500,
		Tone:       types.Tone{Type: types.ToneTypeNone},
		Delay:      2 * time.Second,
		Lockout:    true,
	}, channels[12])
}

func TestSampleChirp(t *testing.T) {
	f := NewChirp()
	channels, err := f.Decode(openSample(t, "chirp.csv"))
	require.NoError(t, err)
	require.Len(t, channels, 12)

	assert.Equal(t, types.Channel{
		Index:      0,
		Name:       "NOAA WX1",
		Frequency:  162400000,
		Modulation: types.ModulationFM,
		Bandwidth:  25000,
		Tone:       types.Tone{Type: types.ToneTypeNone},
		Comment:    "NOAA Weather Radio",
	}, channels[0])

	assert.Equal(t, types.Channel{
		Index:      4,
		Name:       "2m Repeater",
		Frequency:  146940000,
		Modulation: types.ModulationFM,
		Bandwidth:  25000,
		Duplex:     types.DuplexMinus,
		Offset:     600000,
		Tone:       types.Tone{Type: types.ToneTypeCTCSS, Value: 1000},
		Power:      50,
		Comment:    "Local 2m repeater",
	}, channels[4])

	assert.Equal(t, types.Channel{
		Index:      11,
		Name:       "2m DCS RPT",
		Frequency:  147195000,
		Modulation: types.ModulationFM,
		Bandwidth:  25000,
		Duplex:     types.DuplexPlus,
		Offset:     600000,
		Tone:       types.Tone{Type: types.ToneTypeDCS, Value: 71},
		Power:      50,
		Comment:    "2m DCS Repeater",
	}, channels[11])
}

func TestSampleUVPro(t *testing.T) {
	f := NewUVPro()
	channels, err := f.Decode(openSample(t, "uv-pro.csv"))
	require.NoError(t, err)
	require.Len(t, channels, 6)

	assert.Equal(t, types.Channel{
		Index:      1,
		Name:       "VHF Call",
		AlphaTag:   "VHF Call",
		Frequency:  146520000,
		Modulation: types.ModulationFM,
		Bandwidth:  25000,
		Tone:       types.Tone{Type: types.ToneTypeNone},
		Power:      5,
	}, channels[0])

	assert.Equal(t, types.Channel{
		Index:      2,
		Name:       "2m Repeater",
		AlphaTag:   "2m Repeater",
		Frequency:  146940000,
		Modulation: types.ModulationFM,
		Bandwidth:  25000,
		Duplex:     types.DuplexMinus,
		Offset:     600000,
		Tone:       types.Tone{Type: types.ToneTypeCTCSS, Value: 1000},
		Power:      5,
	}, channels[1])

	assert.Equal(t, types.Channel{
		Index:      6,
		Name:       "2m DCS RPT",
		AlphaTag:   "2m DCS RPT",
		Frequency:  147195000,
		Modulation: types.ModulationFM,
		Bandwidth:  25000,
		Duplex:     types.DuplexPlus,
		Offset:     600000,
		Tone:       types.Tone{Type: types.ToneTypeDCS, Value: 71},
		Power:      5,
	}, channels[5])
}

func TestSampleRadioReference(t *testing.T) {
	f := NewRadioReference()
	channels, err := f.Decode(openSample(t, "radioreference.csv"))
	require.NoError(t, err)
	require.Len(t, channels, 14)

	assert.Equal(t, types.Channel{
		Index:      1,
		Name:       "NOAA WX1",
		AlphaTag:   "NOAA1",
		Frequency:  162400000,
		Modulation: types.ModulationFM,
		Bandwidth:  25000,
		Tone:       types.Tone{Type: types.ToneTypeNone},
		Comment:    "Local Government",
	}, channels[0])

	assert.Equal(t, types.Channel{
		Index:      5,
		Name:       "2m Repeater",
		AlphaTag:   "2M RPT",
		Frequency:  146940000,
		Modulation: types.ModulationFM,
		Bandwidth:  25000,
		Duplex:     types.DuplexMinus,
		Offset:     600000,
		Tone:       types.Tone{Type: types.ToneTypeCTCSS, Value: 1000},
		Comment:    "Amateur Radio",
	}, channels[4])

	assert.Equal(t, types.Channel{
		Index:      12,
		Name:       "P25 Dispatch",
		AlphaTag:   "P25DSP",
		Frequency:  851012500,
		Modulation: types.ModulationP25,
		Bandwidth:  12500,
		Tone:       types.Tone{Type: types.ToneTypeNone},
		Comment:    "Public Safety",
	}, channels[11])

	assert.Equal(t, types.Channel{
		Index:      14,
		Name:       "2m DCS Repeater",
		AlphaTag:   "2M DCS",
		Frequency:  147195000,
		Modulation: types.ModulationFM,
		Bandwidth:  25000,
		Duplex:     types.DuplexPlus,
		Offset:     600000,
		Tone:       types.Tone{Type: types.ToneTypeDCS, Value: 71},
		Comment:    "Amateur Radio",
	}, channels[13])
}
