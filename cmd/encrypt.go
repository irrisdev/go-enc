package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/irrisdev/go-enc/genc"
	"github.com/spf13/cobra"
)

var deleteOrigin bool

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a file",
	Long:  `Encrypt a file using AES encryption with the provided passphrase.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(passphrase) < MinPassLen {
			return fmt.Errorf("passphrase must be at least %d characters, got %d", MinPassLen, len(passphrase))
		}

		info, err := os.Stat(file)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("file does not exist: %s", file)
			}
			return fmt.Errorf("error accessing file %s: %w", file, err)
		}

		if info.IsDir() {
			return fmt.Errorf("cannot encrypt directories: %s", file)
		}

		// Check if file is readable
		f, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("file is not readable: %w", err)
		}
		f.Close()

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := genc.Encrypt(passphrase, file, deleteOrigin)

		if deleteOrigin {
			if errors.Is(err, genc.ErrRemoveOrigin) {
				fmt.Println(err)
			} else {
				fmt.Println("original file deleted")

			}
		}

		if err != nil {
			return fmt.Errorf("encryption failed: %w", err)
		}

		fmt.Printf("successfully encrypted: %s\n", file)

		return nil
	},
}

func init() {
	encryptCmd.Flags().BoolVar(&deleteOrigin, "delete-origin", false, "remove original file after encryption")
	rootCmd.AddCommand(encryptCmd)
}
