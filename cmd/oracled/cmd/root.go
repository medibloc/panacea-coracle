package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	homeDir string
	rootCmd = &cobra.Command{
		Use:   "oracled",
		Short: "oracle daemon",
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
	defaultAppHomeDir := filepath.Join(userHomeDir, ".oracle")

	rootCmd.PersistentFlags().StringVar(&homeDir, "home", defaultAppHomeDir, "application home directory")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(startCmd)
}
