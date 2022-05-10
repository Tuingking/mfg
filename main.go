package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/tuingking/mfg/internal/dl"
)

const version = "v0.0.1"

var rooCmd = &cobra.Command{
	Use:     "mfg",
	Short:   "mfg: CLI tool for MFG website.",
	Long:    `mfg: CLI tool for MFG website.`,
	Version: version,
}

func init() {
	rooCmd.AddCommand(dl.CmdDL)
}

func main() {
	if err := rooCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
