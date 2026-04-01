package formats

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/s0lesurviv0r/reband/types"
)

// uvproHeader is the header as written by Go's csv.Writer (no unnecessary quoting).
const uvproHeader = `title,tx_freq,rx_freq,tx_sub_audio(CTCSS=freq/DCS=number),rx_sub_audio(CTCSS=freq/DCS=number),tx_power(H/M/L),bandwidth(12500/25000),scan(0=OFF/1=ON),talk around(0=OFF/1=ON),pre_de_emph_bypass(0=OFF/1=ON),sign(0=OFF/1=ON),tx_dis(0=OFF/1=ON),mute(0=OFF/1=ON),rx_modulation(0=FM/1=AM),tx_modulation(0=FM/1=AM)`

// uvproHeaderQuoted is the header as it appears in UV-Pro files (with quoted "talk around" field).
const uvproHeaderQuoted = `title,tx_freq,rx_freq,tx_sub_audio(CTCSS=freq/DCS=number),rx_sub_audio(CTCSS=freq/DCS=number),tx_power(H/M/L),bandwidth(12500/25000),scan(0=OFF/1=ON),"talk around(0=OFF/1=ON)",pre_de_emph_bypass(0=OFF/1=ON),sign(0=OFF/1=ON),tx_dis(0=OFF/1=ON),mute(0=OFF/1=ON),rx_modulation(0=FM/1=AM),tx_modulation(0=FM/1=AM)`

func TestUVProDecode(t *testing.T) {
	tests := []struct {
		name string
		row  string
		want types.Channel
	}{
		{
			name: "simplex_no_tone",
			row:  `GMRS1,462562500,462562500,0,0,H,25000,1,0,0,1,0,0,0,0`,
			want: types.Channel{
				Index: 1,
				Name:       "GMRS1",
				AlphaTag:   "GMRS1",
				Frequency:  462562500,
				Duplex:     types.DuplexNone,
				Offset:     0,
				Tone:       types.Tone{Type: types.ToneTypeNone},
				Modulation: types.ModulationFM,
				Bandwidth:  25000,
				Power:      5,
				Lockout:    false,
			},
		},
		{
			name: "simplex_ctcss_tone",
			row:  `"ISS RPT UL",145990000,145990000,8850,6700,H,25000,1,0,0,1,0,0,0,0`,
			want: types.Channel{
				Index: 1,
				Name:       "ISS RPT UL",
				AlphaTag:   "ISS RPT UL",
				Frequency:  145990000,
				Duplex:     types.DuplexNone,
				Offset:     0,
				Tone:       types.Tone{Type: types.ToneTypeCTCSS, Value: 670},
				Modulation: types.ModulationFM,
				Bandwidth:  25000,
				Power:      5,
				Lockout:    false,
			},
		},
		{
			name: "positive_offset",
			row:  `"OC RACES",146295000,146895000,8850,13650,H,25000,1,0,0,1,0,0,0,0`,
			want: types.Channel{
				Index: 1,
				Name:       "OC RACES",
				AlphaTag:   "OC RACES",
				Frequency:  146895000,
				Duplex:     types.DuplexMinus,
				Offset:     600000,
				Tone:       types.Tone{Type: types.ToneTypeCTCSS, Value: 1365},
				Modulation: types.ModulationFM,
				Bandwidth:  25000,
				Power:      5,
				Lockout:    false,
			},
		},
		{
			name: "am_modulation",
			row:  `"SNA ATIS",126000000,126000000,0,0,H,25000,1,0,0,1,0,0,1,1`,
			want: types.Channel{
				Index: 1,
				Name:       "SNA ATIS",
				AlphaTag:   "SNA ATIS",
				Frequency:  126000000,
				Duplex:     types.DuplexNone,
				Offset:     0,
				Tone:       types.Tone{Type: types.ToneTypeNone},
				Modulation: types.ModulationAM,
				Bandwidth:  25000,
				Power:      5,
				Lockout:    false,
			},
		},
		{
			name: "locked_out",
			row:  `LOCKED,462562500,462562500,0,0,H,25000,0,0,0,1,0,0,0,0`,
			want: types.Channel{
				Index: 1,
				Name:       "LOCKED",
				AlphaTag:   "LOCKED",
				Frequency:  462562500,
				Duplex:     types.DuplexNone,
				Offset:     0,
				Tone:       types.Tone{Type: types.ToneTypeNone},
				Modulation: types.ModulationFM,
				Bandwidth:  25000,
				Power:      5,
				Lockout:    true,
			},
		},
		{
			name: "ctcss_rounding",
			row:  `"BDF KE6TZG",146985000,146385000,8850,14619,H,25000,1,0,0,1,0,0,0,0`,
			want: types.Channel{
				Index: 1,
				Name:       "BDF KE6TZG",
				AlphaTag:   "BDF KE6TZG",
				Frequency:  146385000,
				Duplex:     types.DuplexPlus,
				Offset:     600000,
				Tone:       types.Tone{Type: types.ToneTypeCTCSS, Value: 1462},
				Modulation: types.ModulationFM,
				Bandwidth:  25000,
				Power:      5,
				Lockout:    false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := uvproHeaderQuoted + "\n" + tt.row + "\n"
			f := NewUVPro()
			channels, err := f.Decode(strings.NewReader(input))
			require.NoError(t, err)
			require.Len(t, channels, 1)
			assert.Equal(t, tt.want, channels[0])
		})
	}
}

func TestUVProEncode(t *testing.T) {
	tests := []struct {
		name    string
		channel types.Channel
		want    string
	}{
		{
			name: "simplex_no_tone",
			channel: types.Channel{
				Name:       "GMRS1",
				Frequency:  462562500,
				Duplex:     types.DuplexNone,
				Tone:       types.Tone{Type: types.ToneTypeNone},
				Modulation: types.ModulationFM,
				Bandwidth:  25000,
				Power:      5,
				Lockout:    false,
			},
			want: uvproHeader + "\nGMRS1,462562500,462562500,0,0,H,25000,1,0,0,1,0,0,0,0\n",
		},
		{
			name: "duplex_with_ctcss",
			channel: types.Channel{
				Name:       "OC RACES",
				Frequency:  146895000,
				Duplex:     types.DuplexMinus,
				Offset:     600000,
				Tone:       types.Tone{Type: types.ToneTypeCTCSS, Value: 1365},
				Modulation: types.ModulationFM,
				Bandwidth:  25000,
				Power:      5,
				Lockout:    false,
			},
			want: uvproHeader + "\nOC RACES,146295000,146895000,13650,13650,H,25000,1,0,0,1,0,0,0,0\n",
		},
		{
			name: "am_locked_out",
			channel: types.Channel{
				Name:       "Air Guard",
				Frequency:  148150000,
				Duplex:     types.DuplexNone,
				Tone:       types.Tone{Type: types.ToneTypeNone},
				Modulation: types.ModulationAM,
				Bandwidth:  25000,
				Power:      2,
				Lockout:    true,
			},
			want: uvproHeader + "\nAir Guard,148150000,148150000,0,0,M,25000,0,0,0,1,0,0,1,1\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewUVPro()
			var buf strings.Builder
			err := f.Encode(&buf, []types.Channel{tt.channel})
			require.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
