package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	var debug bool

	cmd := &cobra.Command{
		Use:   "freq-conv",
		Short: "Convert scanner/transceiver channels between different formats",
		Long:  "Convert scanner/transceiver channels between different formats",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.SetOutput(os.Stdout)
			log.SetFormatter(&log.TextFormatter{
				FullTimestamp: true,
			})
			if debug {
				log.SetLevel(log.DebugLevel)
			}
		},
	}

	cmd.Flags().BoolVar(&debug, "debug", false, "Enable debug logging")

	cmd.AddCommand(DecodeCommand())
	cmd.AddCommand(ConvertCommand())

	return cmd
}
