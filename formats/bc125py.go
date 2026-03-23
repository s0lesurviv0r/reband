package formats

import (
	"fmt"
	"strconv"
	"time"

	"github.com/s0lesurviv0r/reband/types"
)

var bc125pyToStandard = map[string]types.Modulation{
	"nfm": types.ModulationNFM,
	"am":  types.ModulationAM,
	"fm":  types.ModulationFM,
}

var standarToBc125py = map[types.Modulation]string{}

func init() {
	for k, v := range bc125pyToStandard {
		standarToBc125py[v] = k
	}
}

type BC125PY struct {
	GenericCSV
}

func NewBC125PY() *BC125PY {
	return &BC125PY{
		GenericCSV{
			header: []string{
				"Index", "Name", "Frequency (MHz)", "Modulation",
				"CTCSS", "Delay (sec)", "Lockout", "Priority",
			},
			rowDecoder: func(row []string, headerMap map[string]int) (types.Channel, error) {
				index, err := strconv.Atoi(row[headerMap["Index"]])
				if err != nil {
					return types.Channel{}, fmt.Errorf("failed to parse index: %w", err)
				}

				freq, err := types.NewFrequencyFromString(row[headerMap["Frequency (MHz)"]])
				if err != nil {
					return types.Channel{}, fmt.Errorf("failed to parse frequency: %w", err)
				}

				modulation, ok := bc125pyToStandard[row[headerMap["Modulation"]]]
				if !ok {
					return types.Channel{}, ErrUnsupportedModulation
				}

				tone := types.Tone{
					Type:  types.ToneTypeNone,
					Value: 0,
				}

				if row[4] != "none" {
					tone.Type = types.ToneTypeCTCSS
					v, err := strconv.ParseFloat(row[headerMap["CTCSS"]][6:], 64)
					if err != nil {
						return types.Channel{}, fmt.Errorf("failed to parse CTCSS tone: %w", err)
					}
					tone.Value = int(v * 10)
				}

				delay, err := time.ParseDuration(row[headerMap["Delay (sec)"]] + "s")
				if err != nil {
					return types.Channel{}, err
				}

				return types.Channel{
					Index:      index,
					Name:       row[headerMap["Name"]],
					Frequency:  freq,
					Modulation: modulation,
					Tone:       tone,
					Delay:      delay,
					Lockout:    row[headerMap["Lockout"]] == "locked",
					Priority:   row[headerMap["Priority"]] == "on",
				}, nil
			},
			rowEncoder: func(channel types.Channel) ([]string, error) {
				modulation, ok := standarToBc125py[channel.Modulation]
				if !ok {
					return nil, ErrUnsupportedModulation
				}

				ctcss := "none"
				if channel.Tone.Type == types.ToneTypeCTCSS {
					ctcss = "ctcss_" + channel.Tone.CTCSS()
				}

				priority := "off"
				if channel.Priority {
					priority = "on"
				}

				lockout := "unlocked"
				if channel.Lockout {
					lockout = "locked"
				}

				return []string{
					strconv.Itoa(channel.Index),
					channel.Name,
					channel.Frequency.String(),
					modulation,
					ctcss,
					strconv.Itoa(int(channel.Delay.Seconds())),
					lockout,
					priority,
				}, nil
			},
		},
	}
}
