package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/s0lesurviv0r/reband/formats"
)

func ConvertCommand() *cobra.Command {
	var inputPath string
	var outputPath string
	var fromFormat string
	var toFormat string
	var onError string

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

			policy, err := formats.ParseErrorPolicy(onError)
			if err != nil {
				return err
			}

			src, err := formats.Get(fromFormat)
			if err != nil {
				return err
			}
			src.SetErrorPolicy(policy)

			dst, err := formats.Get(toFormat)
			if err != nil {
				return err
			}
			dst.SetErrorPolicy(policy)

			reader, err := os.Open(inputPath)
			if err != nil {
				return err
			}
			defer reader.Close()

			channels, err := src.Decode(reader)
			if err != nil {
				return err
			}

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

	return cmd
}
