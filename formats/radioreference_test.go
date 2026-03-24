package formats

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/s0lesurviv0r/reband/types"
)

const rrHeader = `Frequency Output,Frequency Input,FCC Callsign,Agency/Category,Description,Alpha Tag,PL Tone,Mode,Class Station Code,Tag`

func TestRadioReferenceDecode(t *testing.T) {
	tests := []struct {
		name string
		row  string
		want types.Channel
	}{
		{
			name: "simplex_fmn_no_tone",
			row:  `163.10000,0.00000,,Common Federal Wide Area Common Use,Itinerant Analog,Fed 163.1000 A,,FMN,BM,Federal`,
			want: types.Channel{
				Name:       "Itinerant Analog",
				AlphaTag:   "Fed 163.1000 A",
				Comment:    "Common Federal Wide Area Common Use",
				Frequency:  163100000,
				Duplex:     types.DuplexNone,
				Offset:     0,
				Tone:       types.Tone{Type: types.ToneTypeNone},
				Modulation: types.ModulationFM,
				Bandwidth:  12500,
			},
		},
		{
			name: "simplex_p25",
			row:  `163.10000,0.00000,,Common Federal Wide Area Common Use,Itinerant Digital,Fed 163.1000 D,,P25,BM,Federal`,
			want: types.Channel{
				Name:       "Itinerant Digital",
				AlphaTag:   "Fed 163.1000 D",
				Comment:    "Common Federal Wide Area Common Use",
				Frequency:  163100000,
				Duplex:     types.DuplexNone,
				Offset:     0,
				Tone:       types.Tone{Type: types.ToneTypeNone},
				Modulation: types.ModulationP25,
				Bandwidth:  12500,
			},
		},
		{
			name: "positive_offset_duplex",
			row:  `409.05000,418.05000,,Common Federal Wide Area Common Use,Itinerant Digital,Fed 409.0500 D,,P25,RM,Federal`,
			want: types.Channel{
				Name:       "Itinerant Digital",
				AlphaTag:   "Fed 409.0500 D",
				Comment:    "Common Federal Wide Area Common Use",
				Frequency:  409050000,
				Duplex:     types.DuplexPlus,
				Offset:     9000000,
				Tone:       types.Tone{Type: types.ToneTypeNone},
				Modulation: types.ModulationP25,
				Bandwidth:  12500,
			},
		},
		{
			name: "ctcss_tone",
			row:  `146.94000,147.54000,,Local Repeater,2m Repeater,2M RPT,88.5,FMN,RM,Amateur`,
			want: types.Channel{
				Name:       "2m Repeater",
				AlphaTag:   "2M RPT",
				Comment:    "Local Repeater",
				Frequency:  146940000,
				Duplex:     types.DuplexPlus,
				Offset:     600000,
				Tone:       types.Tone{Type: types.ToneTypeCTCSS, Value: 885},
				Modulation: types.ModulationFM,
				Bandwidth:  12500,
			},
		},
		{
			name: "dcs_tone",
			row:  `462.55000,0.00000,,GMRS,GMRS Ch1,GMRS 1,D023,FM,BM,GMRS`,
			want: types.Channel{
				Name:       "GMRS Ch1",
				AlphaTag:   "GMRS 1",
				Comment:    "GMRS",
				Frequency:  462550000,
				Duplex:     types.DuplexNone,
				Offset:     0,
				Tone:       types.Tone{Type: types.ToneTypeDCS, Value: 23},
				Modulation: types.ModulationFM,
				Bandwidth:  25000,
			},
		},
		{
			name: "negative_offset",
			row:  `147.09000,146.49000,,Local Repeater,2m Repeater,2M RPT,100.0,FMN,RM,Amateur`,
			want: types.Channel{
				Name:       "2m Repeater",
				AlphaTag:   "2M RPT",
				Comment:    "Local Repeater",
				Frequency:  147090000,
				Duplex:     types.DuplexMinus,
				Offset:     600000,
				Tone:       types.Tone{Type: types.ToneTypeCTCSS, Value: 1000},
				Modulation: types.ModulationFM,
				Bandwidth:  12500,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := rrHeader + "\n" + tt.row + "\n"
			f := NewRadioReference()
			channels, err := f.Decode(strings.NewReader(input))
			require.NoError(t, err)
			require.Len(t, channels, 1)
			assert.Equal(t, tt.want, channels[0])
		})
	}
}

func TestRadioReferenceEncode(t *testing.T) {
	tests := []struct {
		name    string
		channel types.Channel
		want    string
	}{
		{
			name: "simplex_fmn",
			channel: types.Channel{
				Name:       "Itinerant Analog",
				AlphaTag:   "Fed 163.1000 A",
				Comment:    "Common Federal Wide Area Common Use",
				Frequency:  163100000,
				Duplex:     types.DuplexNone,
				Tone:       types.Tone{Type: types.ToneTypeNone},
				Modulation: types.ModulationFM,
				Bandwidth:  12500,
			},
			want: rrHeader + "\n163.10000,0.00000,,Common Federal Wide Area Common Use,Itinerant Analog,Fed 163.1000 A,,FMN,BM,\n",
		},
		{
			name: "positive_offset_with_ctcss",
			channel: types.Channel{
				Name:       "2m Repeater",
				AlphaTag:   "2M RPT",
				Comment:    "Local Repeater",
				Frequency:  146940000,
				Duplex:     types.DuplexPlus,
				Offset:     600000,
				Tone:       types.Tone{Type: types.ToneTypeCTCSS, Value: 885},
				Modulation: types.ModulationFM,
				Bandwidth:  12500,
			},
			want: rrHeader + "\n146.94000,147.54000,,Local Repeater,2m Repeater,2M RPT,88.5,FMN,RM,\n",
		},
		{
			name: "p25_positive_offset",
			channel: types.Channel{
				Name:       "Itinerant Digital",
				AlphaTag:   "Fed 409.0500 D",
				Comment:    "Common Federal Wide Area Common Use",
				Frequency:  409050000,
				Duplex:     types.DuplexPlus,
				Offset:     9000000,
				Tone:       types.Tone{Type: types.ToneTypeNone},
				Modulation: types.ModulationP25,
				Bandwidth:  12500,
			},
			want: rrHeader + "\n409.05000,418.05000,,Common Federal Wide Area Common Use,Itinerant Digital,Fed 409.0500 D,,P25,RM,\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewRadioReference()
			var buf strings.Builder
			err := f.Encode(&buf, []types.Channel{tt.channel})
			require.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
