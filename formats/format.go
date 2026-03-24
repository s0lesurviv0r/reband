package formats

import (
	"fmt"
	"io"

	"github.com/s0lesurviv0r/reband/types"
)

type ErrorPolicy int

const (
	ErrorPolicyExit  ErrorPolicy = iota // stop and return the error (default)
	ErrorPolicySkip                     // skip the bad row
	ErrorPolicyEmpty                    // keep the row as an empty placeholder
)

func ParseErrorPolicy(s string) (ErrorPolicy, error) {
	switch s {
	case "exit":
		return ErrorPolicyExit, nil
	case "skip":
		return ErrorPolicySkip, nil
	case "empty":
		return ErrorPolicyEmpty, nil
	default:
		return 0, fmt.Errorf("unknown --on-error value %q: must be exit, skip, or empty", s)
	}
}

type Format interface {
	Decode(io.Reader) ([]types.Channel, error)
	Encode(io.Writer, []types.Channel) error
	SetErrorPolicy(ErrorPolicy)
}

var formatFactories = map[string]func() Format{
	"bc125py": func() Format { return NewBC125PY() },
	"chirp":   func() Format { return NewChirp() },
	"reband":  func() Format { return NewRebandCSV() },
	"uv-pro":  func() Format { return NewUVPro() },
}

func Get(name string) (Format, error) {
	factory, ok := formatFactories[name]
	if !ok {
		return nil, fmt.Errorf("unknown format %s", name)
	}
	return factory(), nil
}
