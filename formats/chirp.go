package formats

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/s0lesurviv0r/reband/types"
)

var chirpToStandard = map[string]types.Modulation{
	"AM": types.ModulationAM,
	"FM": types.ModulationFM,
}

var standardToChirp = map[types.Modulation]string{}

func init() {
	for k, v := range chirpToStandard {
		standardToChirp[v] = k
	}
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

				modulation, ok := chirpToStandard[row[headerMap["Mode"]]]
				if !ok {
					return types.Channel{}, ErrUnsupportedModulation
				}

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

				/*
					RToneFreq:   row[6],
					CToneFreq:   row[7],
					DtcsCode:    row[8],
					DtcsPolarity: row[9],
					RxDtcsCode:  row[10],
					CrossMode:   row[11],
					TStep:       row[13],
					Skip:        row[14],
				*/

				return types.Channel{
					Index:      index,
					Name:       row[1],
					Frequency:  freq,
					Modulation: modulation,
					Duplex:     duplex,
					Offset:     offset,
					Tone:       tone,
					Power:      power,
					Comment:    row[16],
				}, nil
			},
			rowEncoder: func(ch types.Channel) ([]string, error) {
				return nil, nil
			},
		},
	}
}
