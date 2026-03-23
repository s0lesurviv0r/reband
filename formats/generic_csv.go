package formats

import (
	"encoding/csv"
	"io"

	"github.com/s0lesurviv0r/channel-conv/types"
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

	header := records[0]
	headerMap := make(map[string]int, len(header))
	for i, h := range header {
		headerMap[h] = i
	}

	channels := make([]types.Channel, len(records)-1)
	for i, record := range records[1:] { // Skip header
		channel, err := f.rowDecoder(record, headerMap)
		if err != nil {
			return nil, err
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
