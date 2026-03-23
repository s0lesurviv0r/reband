package cmd

import (
	"fmt"
	"os"

	"github.com/rodaine/table"
	"github.com/spf13/cobra"

	"github.com/s0lesurviv0r/reband/formats"
)

func DecodeCommand() *cobra.Command {
	var path string
	var format string

	cmd := &cobra.Command{
		Use:   "decode",
		Short: "Decode a channel list",
		RunE: func(cmd *cobra.Command, args []string) error {
			if path == "" {
				return fmt.Errorf("--path is required")
			}

			if format == "" {
				return fmt.Errorf("--format is required")
			}

			formater, err := formats.Get(format)
			if err != nil {
				return err
			}

			reader, err := os.Open(path)
			if err != nil {
				return err
			}

			channels, err := formater.Decode(reader)
			if err != nil {
				return err
			}

			tbl := table.New("Index", "Name", "AlphaTag", "Frequency", "Modulation", "Duplex", "Offset", "Tone", "Power", "Delay", "Lockout", "Priority", "Comment")
			for _, ch := range channels {
				tbl.AddRow(
					ch.Index,
					ch.Name,
					ch.AlphaTag,
					ch.Frequency,
					ch.Modulation,
					ch.Duplex,
					ch.Offset,
					ch.Tone,
					ch.Power,
					ch.Delay,
					ch.Lockout,
					ch.Priority,
					ch.Comment,
				)
			}
			tbl.Print()

			return nil
		},
	}

	cmd.Flags().StringVar(&path, "path", "", "Path to frequency list")
	cmd.Flags().StringVar(&format, "format", "", "Format to decode from")

	return cmd
}
