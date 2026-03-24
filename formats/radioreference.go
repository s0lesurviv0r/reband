package formats

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/s0lesurviv0r/reband/types"
)

// rrModeInfo maps RadioReference mode strings to modulation and bandwidth.
type rrModeInfo struct {
	modulation types.Modulation
	bandwidth  int
}

var rrModes = map[string]rrModeInfo{
	"FMN":   {types.ModulationFM, 12500},
	"FM":    {types.ModulationFM, 25000},
	"AM":    {types.ModulationAM, 0},
	"P25":   {types.ModulationP25, 12500},
	"DMR":   {types.ModulationDMR, 12500},
	"NXDN":  {types.ModulationNXDN, 12500},
	"DSTAR": {types.ModulationDSTAR, 12500},
	"DStar": {types.ModulationDSTAR, 12500},
	"D-STAR": {types.ModulationDSTAR, 12500},
	"YSF":   {types.ModulationYSF, 12500},
	"TETRA": {types.ModulationTETRA, 12500},
}

type RadioReference struct {
	GenericCSV
}

func NewRadioReference() *RadioReference {
	return &RadioReference{
		GenericCSV{
			autoIndex: true,
			header: []string{
				"Frequency Output",
				"Frequency Input",
				"FCC Callsign",
				"Agency/Category",
				"Description",
				"Alpha Tag",
				"PL Tone",
				"Mode",
				"Class Station Code",
				"Tag",
			},
			rowDecoder: func(row []string, headerMap map[string]int) (types.Channel, error) {
				rxFreq, err := types.NewFrequencyFromString(row[headerMap["Frequency Output"]])
				if err != nil {
					return types.Channel{}, fmt.Errorf("failed to parse frequency output: %w", err)
				}

				txFreqF, err := strconv.ParseFloat(row[headerMap["Frequency Input"]], 64)
				if err != nil {
					return types.Channel{}, fmt.Errorf("failed to parse frequency input: %w", err)
				}
				txFreq := types.Frequency(txFreqF * 1e6)

				duplex := types.DuplexNone
				var offset types.Frequency
				if txFreq > rxFreq {
					duplex = types.DuplexPlus
					offset = txFreq - rxFreq
				} else if txFreq > 0 && txFreq < rxFreq {
					duplex = types.DuplexMinus
					offset = rxFreq - txFreq
				}

				modeStr := row[headerMap["Mode"]]
				info, ok := rrModes[modeStr]
				if !ok {
					return types.Channel{}, fmt.Errorf("unsupported mode %q: %w", modeStr, ErrUnsupportedModulation)
				}

				tone, err := rrDecodeTone(row[headerMap["PL Tone"]])
				if err != nil {
					return types.Channel{}, err
				}

				return types.Channel{
					Name:       row[headerMap["Description"]],
					AlphaTag:   row[headerMap["Alpha Tag"]],
					Comment:    row[headerMap["Agency/Category"]],
					Frequency:  rxFreq,
					Duplex:     duplex,
					Offset:     offset,
					Tone:       tone,
					Modulation: info.modulation,
					Bandwidth:  info.bandwidth,
				}, nil
			},
			rowEncoder: func(ch types.Channel) ([]string, error) {
				modeStr, err := rrEncodeMode(ch.Modulation, ch.Bandwidth)
				if err != nil {
					return nil, err
				}

				rxFreq := fmt.Sprintf("%.5f", float64(ch.Frequency)/1e6)

				txFreqStr := "0.00000"
				classCode := "BM"
				switch ch.Duplex {
				case types.DuplexPlus:
					txFreqStr = fmt.Sprintf("%.5f", float64(ch.Frequency+ch.Offset)/1e6)
					classCode = "RM"
				case types.DuplexMinus:
					txFreqStr = fmt.Sprintf("%.5f", float64(ch.Frequency-ch.Offset)/1e6)
					classCode = "RM"
				}

				return []string{
					rxFreq,
					txFreqStr,
					"", // FCC Callsign
					ch.Comment,
					ch.Name,
					ch.AlphaTag,
					rrEncodeTone(ch.Tone),
					modeStr,
					classCode,
					"", // Tag
				}, nil
			},
		},
	}
}

// rrDecodeTone parses a RadioReference PL Tone field.
// Empty = no tone, decimal number = CTCSS in Hz, D### = DCS code.
func rrDecodeTone(s string) (types.Tone, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return types.Tone{Type: types.ToneTypeNone}, nil
	}
	if strings.HasPrefix(s, "D") {
		code, err := strconv.Atoi(s[1:])
		if err != nil {
			return types.Tone{}, fmt.Errorf("failed to parse DCS tone %q: %w", s, err)
		}
		return types.Tone{Type: types.ToneTypeDCS, Value: code}, nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return types.Tone{}, fmt.Errorf("failed to parse PL tone %q: %w", s, err)
	}
	return types.Tone{Type: types.ToneTypeCTCSS, Value: int(v * 10)}, nil
}

// rrEncodeTone formats a Tone as a RadioReference PL Tone field.
func rrEncodeTone(t types.Tone) string {
	switch t.Type {
	case types.ToneTypeCTCSS:
		return fmt.Sprintf("%.1f", float64(t.Value)/10)
	case types.ToneTypeDCS:
		return fmt.Sprintf("D%03d", t.Value)
	default:
		return ""
	}
}

// rrEncodeMode maps a modulation + bandwidth to a RadioReference mode string.
func rrEncodeMode(m types.Modulation, bandwidth int) (string, error) {
	switch m {
	case types.ModulationFM:
		if bandwidth > 12500 {
			return "FM", nil
		}
		return "FMN", nil
	case types.ModulationAM:
		return "AM", nil
	case types.ModulationP25:
		return "P25", nil
	case types.ModulationDMR:
		return "DMR", nil
	case types.ModulationNXDN:
		return "NXDN", nil
	case types.ModulationDSTAR:
		return "DSTAR", nil
	case types.ModulationYSF, types.ModulationFUSION:
		return "YSF", nil
	case types.ModulationTETRA:
		return "TETRA", nil
	default:
		return "", fmt.Errorf("unsupported modulation %q: %w", m, ErrUnsupportedModulation)
	}
}
