package formats

import (
	"fmt"
	"strconv"
	"time"

	"github.com/s0lesurviv0r/reband/types"
)

// bc125pyBandwidth maps BC125PY modulation strings to channel bandwidth in Hz.
var bc125pyBandwidth = map[string]int{
	"nfm": 12500,
	"fm":  25000,
	"am":  0,
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

				modStr := row[headerMap["Modulation"]]
				bandwidth, ok := bc125pyBandwidth[modStr]
				if !ok {
					return types.Channel{}, fmt.Errorf("unsupported modulation %q: %w", modStr, ErrUnsupportedModulation)
				}

				modulation := types.ModulationFM
				if modStr == "am" {
					modulation = types.ModulationAM
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
					Bandwidth:  bandwidth,
					Tone:       tone,
					Delay:      delay,
					Lockout:    row[headerMap["Lockout"]] == "locked",
					Priority:   row[headerMap["Priority"]] == "on",
				}, nil
			},
			rowEncoder: func(channel types.Channel) ([]string, error) {
				var modStr string
				switch channel.Modulation {
				case types.ModulationAM:
					modStr = "am"
				case types.ModulationFM:
					if channel.Bandwidth > 12500 {
						modStr = "fm"
					} else {
						modStr = "nfm"
					}
				default:
					return nil, fmt.Errorf("unsupported modulation %q: %w", channel.Modulation, ErrUnsupportedModulation)
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
					modStr,
					ctcss,
					strconv.Itoa(int(channel.Delay.Seconds())),
					lockout,
					priority,
				}, nil
			},
		},
	}
}
