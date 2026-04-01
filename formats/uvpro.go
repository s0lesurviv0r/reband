package formats

import (
	"fmt"
	"strconv"

	"github.com/s0lesurviv0r/reband/types"
)

type UVPro struct {
	GenericCSV
}

func NewUVPro() *UVPro {
	return &UVPro{
		GenericCSV{
			autoIndex: true,
			header: []string{
				"title",
				"tx_freq",
				"rx_freq",
				"tx_sub_audio(CTCSS=freq/DCS=number)",
				"rx_sub_audio(CTCSS=freq/DCS=number)",
				"tx_power(H/M/L)",
				"bandwidth(12500/25000)",
				"scan(0=OFF/1=ON)",
				"talk around(0=OFF/1=ON)",
				"pre_de_emph_bypass(0=OFF/1=ON)",
				"sign(0=OFF/1=ON)",
				"tx_dis(0=OFF/1=ON)",
				"mute(0=OFF/1=ON)",
				"rx_modulation(0=FM/1=AM)",
				"tx_modulation(0=FM/1=AM)",
			},
			rowDecoder: func(row []string, headerMap map[string]int) (types.Channel, error) {
				rxFreqHz, err := strconv.ParseInt(row[headerMap["rx_freq"]], 10, 64)
				if err != nil {
					return types.Channel{}, fmt.Errorf("failed to parse rx_freq: %w", err)
				}

				txFreqHz, err := strconv.ParseInt(row[headerMap["tx_freq"]], 10, 64)
				if err != nil {
					return types.Channel{}, fmt.Errorf("failed to parse tx_freq: %w", err)
				}

				duplex := types.DuplexNone
				var offset types.Frequency
				if txFreqHz > rxFreqHz {
					duplex = types.DuplexPlus
					offset = types.Frequency(txFreqHz - rxFreqHz)
				} else if txFreqHz < rxFreqHz {
					duplex = types.DuplexMinus
					offset = types.Frequency(rxFreqHz - txFreqHz)
				}

				rxToneVal, err := strconv.Atoi(row[headerMap["rx_sub_audio(CTCSS=freq/DCS=number)"]])
				if err != nil {
					return types.Channel{}, fmt.Errorf("failed to parse rx_sub_audio: %w", err)
				}

				bandwidth, err := strconv.Atoi(row[headerMap["bandwidth(12500/25000)"]])
				if err != nil {
					return types.Channel{}, fmt.Errorf("failed to parse bandwidth: %w", err)
				}

				scan, err := strconv.Atoi(row[headerMap["scan(0=OFF/1=ON)"]])
				if err != nil {
					return types.Channel{}, fmt.Errorf("failed to parse scan: %w", err)
				}

				rxMod, err := strconv.Atoi(row[headerMap["rx_modulation(0=FM/1=AM)"]])
				if err != nil {
					return types.Channel{}, fmt.Errorf("failed to parse rx_modulation: %w", err)
				}

				modulation := types.ModulationFM
				if rxMod == 1 {
					modulation = types.ModulationAM
				}

				return types.Channel{
					Name:       row[headerMap["title"]],
					AlphaTag:   row[headerMap["title"]],
					Frequency:  types.Frequency(rxFreqHz),
					Duplex:     duplex,
					Offset:     offset,
					Tone:       uvproDecodeTone(rxToneVal),
					Modulation: modulation,
					Bandwidth:  bandwidth,
					Power:      uvproDecodePower(row[headerMap["tx_power(H/M/L)"]]),
					Lockout:    scan == 0,
				}, nil
			},
			rowEncoder: func(ch types.Channel) ([]string, error) {
				if err := uvproValidateFrequency(ch.Frequency); err != nil {
					return nil, err
				}

				rxFreq := int64(ch.Frequency)
				var txFreq int64
				switch ch.Duplex {
				case types.DuplexPlus:
					txFreq = rxFreq + int64(ch.Offset)
				case types.DuplexMinus:
					txFreq = rxFreq - int64(ch.Offset)
				default:
					txFreq = rxFreq
				}

				if err := uvproValidateFrequency(types.Frequency(txFreq)); err != nil {
					return nil, err
				}

				var rxMod, txMod string
				switch ch.Modulation {
				case types.ModulationAM:
					rxMod, txMod = "1", "1"
				case types.ModulationFM:
					rxMod, txMod = "0", "0"
				default:
					return nil, fmt.Errorf("unsupported modulation %q: %w", ch.Modulation, ErrUnsupportedModulation)
				}

				scan := "1"
				if ch.Lockout {
					scan = "0"
				}

				toneStr := uvproEncodeTone(ch.Tone)

				return []string{
					ch.Name,
					strconv.FormatInt(txFreq, 10),
					strconv.FormatInt(rxFreq, 10),
					toneStr,
					toneStr,
					uvproEncodePower(ch.Power),
					strconv.Itoa(ch.Bandwidth),
					scan,
					"0", // talk around
					"0", // pre_de_emph_bypass
					"1", // sign
					"0", // tx_dis
					"0", // mute
					rxMod,
					txMod,
				}, nil
			},
		},
	}
}

// uvproDecodeTone converts a UV-Pro sub_audio value to a Tone.
// UV-Pro stores CTCSS as frequency × 100 (e.g. 8850 = 88.5 Hz).
// Reband stores CTCSS as frequency × 10 (e.g. 885 = 88.5 Hz).
// Values ≥ 6700 are treated as CTCSS; values 1–599 are treated as DCS codes; 0 = none.
func uvproDecodeTone(v int) types.Tone {
	if v == 0 {
		return types.Tone{Type: types.ToneTypeNone}
	}
	if v >= 6700 {
		// CTCSS: convert from 1/100 Hz to 1/10 Hz (round to nearest)
		return types.Tone{Type: types.ToneTypeCTCSS, Value: (v + 5) / 10}
	}
	// DCS code stored directly
	return types.Tone{Type: types.ToneTypeDCS, Value: v}
}

// uvproEncodeTone converts a Tone to a UV-Pro sub_audio value string.
// CTCSS is stored as frequency × 100; DCS is stored as the code number.
func uvproEncodeTone(t types.Tone) string {
	switch t.Type {
	case types.ToneTypeCTCSS:
		// Convert from 1/10 Hz (reband) to 1/100 Hz (UV-Pro)
		return strconv.Itoa(t.Value * 10)
	case types.ToneTypeDCS:
		return strconv.Itoa(t.Value)
	default:
		return "0"
	}
}

// uvproDecodePower maps UV-Pro power string to approximate watts.
func uvproDecodePower(s string) int {
	switch s {
	case "H":
		return 5
	case "M":
		return 2
	case "L":
		return 1
	default:
		return 0
	}
}

// uvproValidateFrequency checks that a frequency is within UV-Pro supported ranges:
// VHF 136–174 MHz or UHF 400–520 MHz.
func uvproValidateFrequency(f types.Frequency) error {
	mhz := float64(f) / 1e6
	if (mhz >= 136 && mhz <= 174) || (mhz >= 400 && mhz <= 520) {
		return nil
	}
	return fmt.Errorf("frequency %.4f MHz is outside UV-Pro supported ranges (VHF: 136–174 MHz, UHF: 400–520 MHz)", mhz)
}

// uvproEncodePower maps watts to UV-Pro power level.
func uvproEncodePower(w int) string {
	switch {
	case w >= 5:
		return "H"
	case w >= 2:
		return "M"
	case w >= 1:
		return "L"
	default:
		return "H"
	}
}
