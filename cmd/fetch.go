package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/motextur3/dom-distiller/distiller"
	"github.com/motextur3/dom-distiller/fetcher"
	"github.com/motextur3/dom-distiller/formatter"
)

var (
	format  string
	waitFor string
	timeout int
)

var fetchCmd = &cobra.Command{
	Use:   "fetch [url]",
	Short: "Fetch and distill a URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		ctxTimeout := time.Duration(timeout) * time.Second

		// 1. Fetch
		html, err := fetcher.Fetch(url, waitFor, ctxTimeout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching URL: %v\n", err)
			os.Exit(1)
		}

		// 2. Distill
		node, err := distiller.Distill(html)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error distilling HTML: %v\n", err)
			os.Exit(1)
		}

		// 3. Format
		var output string
		switch format {
		case "markdown":
			output = formatter.ToMarkdown(node)
		case "json":
			output = formatter.ToJSON(node)
		default:
			fmt.Fprintf(os.Stderr, "Unknown format: %s\n", format)
			os.Exit(1)
		}

		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().StringVarP(&format, "format", "f", "markdown", "Output format (markdown|json)")
	fetchCmd.Flags().StringVarP(&waitFor, "wait-for", "w", "", "CSS selector to wait for before extracting")
	fetchCmd.Flags().IntVarP(&timeout, "timeout", "t", 30, "Timeout in seconds")
}