package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dom-distiller",
	Short: "A CLI tool to distill web pages into token-optimized Agent Views.",
	Long: `dom-distiller takes a URL, renders it using a headless browser,
and outputs a highly compressed, semantic view of the page, stripping out noise
and highlighting actionable elements for AI agents.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}