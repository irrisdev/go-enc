package cmd

import (
	"errors"
	"os"

	"github.com/irrisdev/go-enc/genc"
	"github.com/spf13/cobra"
)

var deleteOrigin bool

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a file",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(passphrase) < MinPassLen {
			return errors.New("passphrase too short")
		}

		info, err := os.Stat(file)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return errors.New("directories are not supported")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		genc.Encrypt(passphrase, file, deleteOrigin)
	},
}

func init() {
	encryptCmd.Flags().BoolVar(&deleteOrigin, "delete-origin", false, "remove original file after encryption")
	rootCmd.AddCommand(encryptCmd)
}
