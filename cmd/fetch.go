package cmd

import (
	"encoding/json"
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
	saveMap string
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

		// 3. Save Map (Optional)
		if saveMap != "" {
			actionMap := make(map[string]string)
			buildActionMap(node, actionMap)

			mapData, err := json.MarshalIndent(actionMap, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating action map: %v\n", err)
			} else {
				err = os.WriteFile(saveMap, mapData, 0644)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error writing action map to %s: %v\n", saveMap, err)
				}
			}
		}

		// 4. Format
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

func buildActionMap(n *distiller.Node, m map[string]string) {
	if n.ActionID != "" && n.XPath != "" {
		m[n.ActionID] = n.XPath
	}
	for _, c := range n.Children {
		buildActionMap(c, m)
	}
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().StringVarP(&format, "format", "f", "markdown", "Output format (markdown|json)")
	fetchCmd.Flags().StringVarP(&waitFor, "wait-for", "w", "", "CSS selector to wait for before extracting")
	fetchCmd.Flags().IntVarP(&timeout, "timeout", "t", 30, "Timeout in seconds")
	fetchCmd.Flags().StringVarP(&saveMap, "save-map", "s", "", "Path to save the ActionID -> XPath JSON map (e.g., /tmp/map.json)")
}