package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const MinPassLen = 10

var (
	passphrase string
	file       string
)

var rootCmd = &cobra.Command{
	Use:   "go-enc",
	Short: "Encrypt and decrypt files",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "path to file")
	rootCmd.PersistentFlags().StringVarP(&passphrase, "passphrase", "p", "", "encryption passphrase")

	rootCmd.MarkPersistentFlagRequired("file")
	rootCmd.MarkPersistentFlagRequired("passphrase")
}
