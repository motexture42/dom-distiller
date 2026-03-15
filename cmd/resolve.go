package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var resolveCmd = &cobra.Command{
	Use:   "resolve [map_file] [action_id]",
	Short: "Resolve an ActionID to its XPath using a saved map file",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		mapFile := args[0]
		actionID := args[1]

		data, err := os.ReadFile(mapFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading map file: %v\n", err)
			os.Exit(1)
		}

		var actionMap map[string]string
		if err := json.Unmarshal(data, &actionMap); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing map file: %v\n", err)
			os.Exit(1)
		}

		xpath, exists := actionMap[actionID]
		if !exists {
			fmt.Fprintf(os.Stderr, "Error: ActionID '%s' not found in map.\n", actionID)
			os.Exit(1)
		}

		// Print only the XPath to stdout so other tools can pipe it
		fmt.Print(xpath)
	},
}

func init() {
	rootCmd.AddCommand(resolveCmd)
}