package formats

import (
	"fmt"
	"strconv"
	"time"

	"github.com/s0lesurviv0r/reband/types"
)

var modulationToString = map[types.Modulation]string{
	types.ModulationFM:     "fm",
	types.ModulationNFM:    "nfm",
	types.ModulationAM:     "am",
	types.ModulationWFM:    "wfm",
	types.ModulationLSB:    "lsb",
	types.ModulationUSB:    "usb",
	types.ModulationCW:     "cw",
	types.ModulationC4FM:   "c4fm",
	types.ModulationDSTAR:  "dstar",
	types.ModulationP25:    "p25",
	types.ModulationNXDN:   "nxdn",
	types.ModulationDMR:    "dmr",
	types.ModulationYSF:    "ysf",
	types.ModulationFUSION: "fusion",
	types.ModulationPOCSAG: "pocsag",
	types.ModulationDPMR:   "dpmr",
	types.ModulationTETRA:  "tetra",
}

var stringToModulation = map[string]types.Modulation{}

func init() {
	for k, v := range modulationToString {
		stringToModulation[v] = k
	}
}

type RebandCSV struct {
	GenericCSV
}

func NewRebandCSV() *RebandCSV {
	return &RebandCSV{
		GenericCSV{
			header: []string{
				"Index", "Name", "AlphaTag", "Comment",
				"Frequency", "Duplex", "Offset",
				"ToneType", "ToneValue",
				"Modulation", "Power", "Delay", "Lockout", "Priority",
			},
			rowDecoder: func(row []string, headerMap map[string]int) (types.Channel, error) {
				index, err := strconv.Atoi(row[headerMap["Index"]])
				if err != nil {
					return types.Channel{}, fmt.Errorf("failed to parse index: %w", err)
				}

				freq, err := types.NewFrequencyFromString(row[headerMap["Frequency"]])
				if err != nil {
					return types.Channel{}, fmt.Errorf("failed to parse frequency: %w", err)
				}

				duplex := types.DuplexNone
				switch row[headerMap["Duplex"]] {
				case "+":
					duplex = types.DuplexPlus
				case "-":
					duplex = types.DuplexMinus
				}

				offset, err := types.NewFrequencyFromString(row[headerMap["Offset"]])
				if err != nil {
					return types.Channel{}, fmt.Errorf("failed to parse offset: %w", err)
				}

				tone := types.Tone{Type: types.ToneTypeNone}
				switch row[headerMap["ToneType"]] {
				case "ctcss":
					v, err := strconv.Atoi(row[headerMap["ToneValue"]])
					if err != nil {
						return types.Channel{}, fmt.Errorf("failed to parse tone value: %w", err)
					}
					tone = types.Tone{Type: types.ToneTypeCTCSS, Value: v}
				case "dcs":
					v, err := strconv.Atoi(row[headerMap["ToneValue"]])
					if err != nil {
						return types.Channel{}, fmt.Errorf("failed to parse tone value: %w", err)
					}
					tone = types.Tone{Type: types.ToneTypeDCS, Value: v}
				}

				modulation, ok := stringToModulation[row[headerMap["Modulation"]]]
				if !ok {
					return types.Channel{}, ErrUnsupportedModulation
				}

				power, err := strconv.Atoi(row[headerMap["Power"]])
				if err != nil {
					return types.Channel{}, fmt.Errorf("failed to parse power: %w", err)
				}

				delay, err := time.ParseDuration(row[headerMap["Delay"]] + "s")
				if err != nil {
					return types.Channel{}, fmt.Errorf("failed to parse delay: %w", err)
				}

				return types.Channel{
					Index:      index,
					Name:       row[headerMap["Name"]],
					AlphaTag:   row[headerMap["AlphaTag"]],
					Comment:    row[headerMap["Comment"]],
					Frequency:  freq,
					Duplex:     duplex,
					Offset:     offset,
					Tone:       tone,
					Modulation: modulation,
					Power:      power,
					Delay:      delay,
					Lockout:    row[headerMap["Lockout"]] == "true",
					Priority:   row[headerMap["Priority"]] == "true",
				}, nil
			},
			rowEncoder: func(ch types.Channel) ([]string, error) {
				modStr, ok := modulationToString[ch.Modulation]
				if !ok {
					return nil, ErrUnsupportedModulation
				}

				duplexStr := ""
				switch ch.Duplex {
				case types.DuplexPlus:
					duplexStr = "+"
				case types.DuplexMinus:
					duplexStr = "-"
				}

				toneType := "none"
				toneValue := "0"
				switch ch.Tone.Type {
				case types.ToneTypeCTCSS:
					toneType = "ctcss"
					toneValue = strconv.Itoa(ch.Tone.Value)
				case types.ToneTypeDCS:
					toneType = "dcs"
					toneValue = strconv.Itoa(ch.Tone.Value)
				}

				lockout := "false"
				if ch.Lockout {
					lockout = "true"
				}

				priority := "false"
				if ch.Priority {
					priority = "true"
				}

				return []string{
					strconv.Itoa(ch.Index),
					ch.Name,
					ch.AlphaTag,
					ch.Comment,
					ch.Frequency.String(),
					duplexStr,
					ch.Offset.String(),
					toneType,
					toneValue,
					modStr,
					strconv.Itoa(ch.Power),
					strconv.Itoa(int(ch.Delay.Seconds())),
					lockout,
					priority,
				}, nil
			},
		},
	}
}
