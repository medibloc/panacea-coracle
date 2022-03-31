package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	homeDir string
	rootCmd = &cobra.Command{
		Use:   "datavald",
		Short: "panacea-data-market-validator",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

// init is run automatically when the package is loaded.
func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	defaultAppHomeDir := filepath.Join(userHomeDir, ".dataval")

	rootCmd.PersistentFlags().StringVar(&homeDir, "home", defaultAppHomeDir, "application home directory")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(startCmd)
}
