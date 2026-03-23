package formats

import (
	"fmt"
	"io"

	"github.com/s0lesurviv0r/reband/types"
)

type Format interface {
	Decode(io.Reader) ([]types.Channel, error)
	Encode(io.Writer, []types.Channel) error
}

var formats = map[string]Format{
	"bc125py":   NewBC125PY(),
	"chirp":     NewChirp(),
	"reband": NewRebandCSV(),
}

func Get(name string) (Format, error) {
	format, ok := formats[name]
	if !ok {
		return nil, fmt.Errorf("unknown format %s", name)
	}

	return format, nil
}
