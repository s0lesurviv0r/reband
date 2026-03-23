package types

import (
	"fmt"
	"strconv"
	"time"
)

// Frequency is a type that represents a frequency in Hz
type Frequency int64

func NewFrequencyFromString(f string) (Frequency, error) {
	v, err := strconv.ParseFloat(f, 64)
	if err != nil {
		return 0, err
	}

	return Frequency(v * 1e6), nil
}

func (f Frequency) String() string {
	return f.StringWithPrecision(4)
}

func (f Frequency) StringWithPrecision(precision int) string {
	format := "%." + strconv.Itoa(precision) + "f"
	return fmt.Sprintf(format, float64(f)/1e6)
}

type ToneType int

const (
	ToneTypeNone ToneType = iota
	ToneTypeCTCSS
	ToneTypeDCS
)

type Duplex int

const (
	DuplexNone Duplex = iota
	DuplexPlus
	DuplexMinus
)

func (d Duplex) String() string {
	switch d {
	case DuplexPlus:
		return "+"
	case DuplexMinus:
		return "-"
	default:
		return ""
	}
}

type Tone struct {
	Type ToneType

	// Value is the CTCSS or DCS value. For
	// CTCSS, it is the frequency in 1/10 Hz. For DCS,
	// it is the DCS code
	Value int
}

func (t Tone) CTCSS() string {
	return fmt.Sprintf("%.1f", float64(t.Value)/10)
}

func (t Tone) String() string {
	switch t.Type {
	case ToneTypeCTCSS:
		return "CTCSS " + t.CTCSS()
	case ToneTypeDCS:
		return fmt.Sprintf("DCS %03d", t.Value)
	default:
		return ""
	}
}

type Channel struct {
	Index      int
	Name       string
	AlphaTag   string
	Comment    string
	Frequency  Frequency
	Duplex     Duplex
	Offset     Frequency
	Tone       Tone
	Modulation Modulation
	Power      int
	Delay      time.Duration
	Lockout    bool
	Priority   bool
}
