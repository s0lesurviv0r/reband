package formats

import (
	"fmt"
	"io"

	"github.com/s0lesurviv0r/channel-conv/types"
)

type Format interface {
	Decode(io.Reader) ([]types.Channel, error)
	Encode(io.Writer, []types.Channel) error
}

var formats = map[string]Format{
	"bc125py": &BC125PY{},
	"chirp":   &Chirp{},
}

func Get(name string) (Format, error) {
	format, ok := formats[name]
	if !ok {
		return nil, fmt.Errorf("unknown format %s", name)
	}

	return format, nil
}
