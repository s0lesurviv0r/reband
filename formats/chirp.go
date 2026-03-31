package formats

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/s0lesurviv0r/reband/types"
)

// chirpModeInfo holds the modulation and bandwidth for a CHIRP mode string.
type chirpModeInfo struct {
	modulation types.Modulation
	bandwidth  int
}

// chirpModes maps CHIRP mode strings to modulation and bandwidth.
var chirpModes = map[string]chirpModeInfo{
	"FM":  {types.ModulationFM, 25000},
	"NFM": {types.ModulationFM, 12500},
	"AM":  {types.ModulationAM, 0},
	"USB": {types.ModulationUSB, 0},
	"LSB": {types.ModulationLSB, 0},
	"CW":  {types.ModulationCW, 0},
	"CWR": {types.ModulationCW, 0},
	"DV":  {types.ModulationDSTAR, 0},
	"WFM": {types.ModulationWFM, 0},
	"P25": {types.ModulationP25, 12500},
}

type Chirp struct {
	GenericCSV
}

func NewChirp() *Chirp {
	return &Chirp{
		GenericCSV{
			header: []string{
				"Location",
				"Name",
				"Frequency",
				"Duplex",
				"Offset",
				"Tone",
				"rToneFreq",
				"cToneFreq",
				"DtcsCode",
				"DtcsPolarity",
				"RxDtcsCode",
				"CrossMode",
				"Mode",
				"TStep",
				"Skip",
				"Power",
				"Comment",
				"URCALL",
				"RPT1CALL",
				"RPT2CALL",
				"DVCODE",
			},
			rowDecoder: func(row []string, headerMap map[string]int) (types.Channel, error) {
				if len(row) < 21 {
					return types.Channel{}, fmt.Errorf("row has insufficient columns")
				}

				index, err := strconv.Atoi(row[headerMap["Location"]])
				if err != nil {
					return types.Channel{}, fmt.Errorf("failed to parse index: %w", err)
				}

				freq, err := types.NewFrequencyFromString(row[headerMap["Frequency"]])
				if err != nil {
					return types.Channel{}, fmt.Errorf("failed to parse frequency: %w", err)
				}

				modeStr := row[headerMap["Mode"]]
				modeInfo, ok := chirpModes[modeStr]
				if !ok {
					return types.Channel{}, fmt.Errorf("unsupported modulation %q: %w", modeStr, ErrUnsupportedModulation)
				}
				modulation := modeInfo.modulation
				bandwidth := modeInfo.bandwidth

				tone := types.Tone{
					Type:  types.ToneTypeNone,
					Value: 0,
				}

				toneType := row[headerMap["Tone"]]
				switch toneType {
				case "Tone":
					tone.Type = types.ToneTypeCTCSS
					v, err := strconv.ParseFloat(row[headerMap["rToneFreq"]], 64)
					if err != nil {
						return types.Channel{}, fmt.Errorf("failed to parse CTCSS tone: %w", err)
					}
					tone.Value = int(v * 10)
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

				power, err := strconv.Atoi(strings.TrimSuffix(row[headerMap["Power"]], "W"))
				if err != nil {
					return types.Channel{}, fmt.Errorf("failed to parse power: %w", err)
				}

				return types.Channel{
					Index:      index,
					Name:       row[1],
					Frequency:  freq,
					Modulation: modulation,
					Bandwidth:  bandwidth,
					Duplex:     duplex,
					Offset:     offset,
					Tone:       tone,
					Power:      power,
					Comment:    row[16],
				}, nil
			},
			rowEncoder: func(ch types.Channel) ([]string, error) {
				var mode string
				switch ch.Modulation {
				case types.ModulationAM:
					mode = "AM"
				case types.ModulationFM:
					if ch.Bandwidth > 12500 {
						mode = "FM"
					} else {
						mode = "NFM"
					}
				case types.ModulationUSB:
					mode = "USB"
				case types.ModulationLSB:
					mode = "LSB"
				case types.ModulationCW:
					mode = "CW"
				case types.ModulationDSTAR:
					mode = "DV"
				case types.ModulationWFM:
					mode = "WFM"
				case types.ModulationP25:
					mode = "P25"
				default:
					return nil, fmt.Errorf("unsupported modulation %q: %w", ch.Modulation, ErrUnsupportedModulation)
				}

				duplexStr := ""
				switch ch.Duplex {
				case types.DuplexPlus:
					duplexStr = "+"
				case types.DuplexMinus:
					duplexStr = "-"
				}

				toneField := ""
				rToneFreq := "88.5"
				cToneFreq := "88.5"
				dtcsCode := "023"
				dtcsPolarity := "NN"
				rxDtcsCode := "023"

				switch ch.Tone.Type {
				case types.ToneTypeCTCSS:
					toneField = "Tone"
					rToneFreq = ch.Tone.CTCSS()
				case types.ToneTypeDCS:
					toneField = "DTCS"
					dtcsCode = fmt.Sprintf("%03d", ch.Tone.Value)
					rxDtcsCode = dtcsCode
				}

				return []string{
					strconv.Itoa(ch.Index),
					ch.Name,
					fmt.Sprintf("%f", float64(ch.Frequency)/1e6),
					duplexStr,
					fmt.Sprintf("%f", float64(ch.Offset)/1e6),
					toneField,
					rToneFreq,
					cToneFreq,
					dtcsCode,
					dtcsPolarity,
					rxDtcsCode,
					"Tone->Tone",
					mode,
					"5.00",
					"",
					fmt.Sprintf("%dW", ch.Power),
					ch.Comment,
					"",
					"",
					"",
					"",
				}, nil
			},
		},
	}
}
