package formats

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/s0lesurviv0r/reband/types"
)

type GenericCSV struct {
	header      []string
	rowEncoder  func(channel types.Channel) ([]string, error)
	rowDecoder  func(record []string, headerMap map[string]int) (types.Channel, error)
	errorPolicy ErrorPolicy
	autoIndex   bool // assign 1-based Index after decoding when format has no index column
}

func (f *GenericCSV) SetErrorPolicy(p ErrorPolicy) {
	f.errorPolicy = p
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

	channels := make([]types.Channel, 0, len(records)-1)
	for i, record := range records[1:] { // Skip header
		rowNum := i + 1

		if len(record) < len(f.header) {
			rowErr := fmt.Errorf("row %d has %d columns, expected %d", rowNum, len(record), len(f.header))
			switch f.errorPolicy {
			case ErrorPolicySkip:
				fmt.Fprintf(os.Stderr, "warning: skipping row %d: %v\n", rowNum, rowErr)
				continue
			case ErrorPolicyEmpty:
				fmt.Fprintf(os.Stderr, "warning: row %d set to empty: %v\n", rowNum, rowErr)
				channels = append(channels, types.Channel{})
				continue
			default:
				return nil, rowErr
			}
		}

		channel, err := f.rowDecoder(record, headerMap)
		if err != nil {
			rowErr := fmt.Errorf("row %d: %w", rowNum, err)
			switch f.errorPolicy {
			case ErrorPolicySkip:
				fmt.Fprintf(os.Stderr, "warning: skipping row %d: %v\n", rowNum, err)
				continue
			case ErrorPolicyEmpty:
				fmt.Fprintf(os.Stderr, "warning: row %d set to empty: %v\n", rowNum, err)
				channels = append(channels, types.Channel{})
				continue
			default:
				return nil, rowErr
			}
		}

		channels = append(channels, channel)
	}

	if f.autoIndex {
		for i := range channels {
			channels[i].Index = i + 1
		}
	}

	return channels, nil
}

func (f *GenericCSV) Encode(writer io.Writer, channels []types.Channel) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	if err := csvWriter.Write(f.header); err != nil {
		return err
	}

	emptyRow := make([]string, len(f.header))

	for i, channel := range channels {
		record, err := f.rowEncoder(channel)
		if err != nil {
			switch f.errorPolicy {
			case ErrorPolicySkip:
				fmt.Fprintf(os.Stderr, "warning: skipping channel %d: %v\n", i+1, err)
				continue
			case ErrorPolicyEmpty:
				fmt.Fprintf(os.Stderr, "warning: channel %d set to empty: %v\n", i+1, err)
				record = emptyRow
			default:
				return err
			}
		}

		if err := csvWriter.Write(record); err != nil {
			return err
		}
	}

	return nil
}
