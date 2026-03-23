package formats

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/s0lesurviv0r/reband/types"
)

type GenericCSV struct {
	header     []string
	rowEncoder func(channel types.Channel) ([]string, error)
	rowDecoder func(record []string, headerMap map[string]int) (types.Channel, error)
}

func (f *GenericCSV) Decode(reader io.Reader) ([]types.Channel, error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("empty file")
	}

	header := records[0]
	headerMap := make(map[string]int, len(header))
	for i, h := range header {
		headerMap[h] = i
	}

	for _, col := range f.header {
		if _, ok := headerMap[col]; !ok {
			return nil, fmt.Errorf("missing required column %q", col)
		}
	}

	channels := make([]types.Channel, len(records)-1)
	for i, record := range records[1:] { // Skip header
		if len(record) < len(f.header) {
			return nil, fmt.Errorf("row %d has %d columns, expected %d", i+1, len(record), len(f.header))
		}
		channel, err := f.rowDecoder(record, headerMap)
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", i+1, err)
		}
		channels[i] = channel
	}

	return channels, nil
}

func (f *GenericCSV) Encode(writer io.Writer, channels []types.Channel) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	if err := csvWriter.Write(f.header); err != nil {
		return err
	}

	for _, channel := range channels {
		record, err := f.rowEncoder(channel)
		if err != nil {
			return err
		}

		err = csvWriter.Write(record)
		if err != nil {
			return err
		}
	}

	return nil
}
