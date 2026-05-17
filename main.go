package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

const (
	// Version is the current version of vhs.
	Version = "0.1.0"

	// DefaultTapeFile is the default tape file extension.
	DefaultTapeFile = ".tape"

	// DefaultOutputFile is the default output file name.
	// Changed from output.gif to recording.gif to be more descriptive.
	DefaultOutputFile = "recording.gif"
)

var rootCmd = &cobra.Command{
	Use:          "vhs <file>",
	Short:        "Run a VHS tape file to generate a terminal GIF",
	Version:      Version,
	Args:         cobra.MaximumNArgs(1),
	RunE:         run,
	SilenceUsage: true,
}

func init() {
	rootCmd.Flags().StringP("output", "o", DefaultOutputFile, "output file path")
	rootCmd.Flags().BoolP("publish", "p", false, "publish the GIF to vhs.charm.sh")
	rootCmd.Flags().BoolP("quiet", "q", false, "quiet mode (suppress output)")
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("no tape file provided")
	}

	tapeFile := args[0]

	// Validate the tape file exists
	if _, err := os.Stat(tapeFile); os.IsNotExist(err) {
		return fmt.Errorf("tape file not found: %s", tapeFile)
	}

	quiet, _ := cmd.Flags().GetBool("quiet")
	if !quiet {
		log.Info("Processing tape", "file", tapeFile)
	}

	// Read the tape file
	data, err := os.ReadFile(tapeFile)
	if err != nil {
		return fmt.Errorf("failed to read tape file: %w", err)
	}

	output, _ := cmd.Flags().GetString("output")

	_ = data

	if !quiet {
		log.Info("Done!", "output", output)
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Error("Error", "err", err)
		os.Exit(1)
	}
}
