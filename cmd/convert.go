package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/s0lesurviv0r/reband/formats"
	"github.com/s0lesurviv0r/reband/types"
)

func ConvertCommand() *cobra.Command {
	var inputPath string
	var outputPath string
	var fromFormat string
	var toFormat string
	var onError string
	var splitSize int
	var outputDir string

	cmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert a channel list between formats",
		RunE: func(cmd *cobra.Command, args []string) error {
			if inputPath == "" {
				_ = cmd.Usage()
				return fmt.Errorf("--input is required")
			}
			if fromFormat == "" {
				_ = cmd.Usage()
				return fmt.Errorf("--from is required")
			}
			if toFormat == "" {
				_ = cmd.Usage()
				return fmt.Errorf("--to is required")
			}
			if splitSize > 0 && outputDir == "" {
				_ = cmd.Usage()
				return fmt.Errorf("--output-dir is required when --split-output-size is set")
			}
			if splitSize > 0 && outputPath != "" {
				_ = cmd.Usage()
				return fmt.Errorf("--output and --split-output-size are mutually exclusive")
			}

			policy, err := formats.ParseErrorPolicy(onError)
			if err != nil {
				return err
			}

			src, err := formats.Get(fromFormat)
			if err != nil {
				return err
			}
			src.SetErrorPolicy(policy)

			reader, err := os.Open(inputPath)
			if err != nil {
				return err
			}
			defer reader.Close()

			channels, err := src.Decode(reader)
			if err != nil {
				return err
			}

			if splitSize > 0 {
				return encodeChunks(channels, splitSize, outputDir, toFormat, policy)
			}

			dst, err := formats.Get(toFormat)
			if err != nil {
				return err
			}
			dst.SetErrorPolicy(policy)

			writer := os.Stdout
			if outputPath != "" {
				f, err := os.Create(outputPath)
				if err != nil {
					return err
				}
				defer f.Close()
				writer = f
			}

			return dst.Encode(writer, channels)
		},
	}

	cmd.Flags().StringVar(&inputPath, "input", "", "Path to input file")
	cmd.Flags().StringVar(&outputPath, "output", "", "Path to output file (defaults to stdout)")
	cmd.Flags().StringVar(&fromFormat, "from", "", "Source format")
	cmd.Flags().StringVar(&toFormat, "to", "", "Destination format")
	cmd.Flags().StringVar(&onError, "on-error", "exit", "How to handle row errors: exit, skip, or empty")
	cmd.Flags().IntVar(&splitSize, "split-output-size", 0, "Split output into multiple files with at most this many channels each")
	cmd.Flags().StringVar(&outputDir, "output-dir", "", "Directory to write split output files (required with --split-output-size)")

	return cmd
}

func encodeChunks(channels []types.Channel, size int, dir string, toFormat string, policy formats.ErrorPolicy) error {
	// Pre-filter channels through the encoder before splitting so that chunk
	// sizes reflect the valid channel count, not the raw decoded count.
	dst, err := formats.Get(toFormat)
	if err != nil {
		return err
	}
	dst.SetErrorPolicy(policy)
	channels, err = dst.FilterChannels(channels)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("failed to clear output directory: %w", err)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	numFiles := (len(channels) + size - 1) / size
	digits := len(strconv.Itoa(numFiles))

	for i := range numFiles {
		chunk := channels[i*size : min(i*size+size, len(channels))]
		filename := fmt.Sprintf("%0*d.csv", digits, i+1)
		path := filepath.Join(dir, filename)

		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create %s: %w", path, err)
		}

		dst, err := formats.Get(toFormat)
		if err != nil {
			f.Close()
			return err
		}
		dst.SetErrorPolicy(policy)

		if err := dst.Encode(f, chunk); err != nil {
			f.Close()
			return fmt.Errorf("failed to encode %s: %w", path, err)
		}

		if err := f.Close(); err != nil {
			return fmt.Errorf("failed to close %s: %w", path, err)
		}
	}

	return nil
}

