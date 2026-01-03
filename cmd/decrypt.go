package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/irrisdev/go-enc/genc"
	"github.com/spf13/cobra"
)

var outPath string

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt a file",
	Long:  `Decrypt a file that was encrypted with this tool.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Check if input file exists
		info, err := os.Stat(file)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("file does not exist: %s", file)
			}
			return fmt.Errorf("error accessing file %s: %w", file, err)
		}

		if info.IsDir() {
			return fmt.Errorf("cannot decrypt directories: %s", file)
		}

		// Validate output path if provided
		if outPath != "" {
			// Check if output directory exists
			outDir := filepath.Dir(outPath)
			if outDir != "." && outDir != "" {
				if _, err := os.Stat(outDir); err != nil {
					if os.IsNotExist(err) {
						return fmt.Errorf("output directory does not exist: %s", outDir)
					}
					return fmt.Errorf("error accessing output directory: %w", err)
				}
			}

			// Check if output file already exists
			if _, err := os.Stat(outPath); err == nil {
				return fmt.Errorf("output file already exists: %s", outPath)
			}
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		if outPath == "" {
			err = genc.Decrypt(passphrase, file)
		} else {
			err = genc.Decrypt(passphrase, file, outPath)
		}

		if err != nil {
			return fmt.Errorf("decryption failed: %w", err)
		}

		if outPath == "" {
			fmt.Printf("successfully decrypted: %s\n", file)
		} else {
			fmt.Printf("successfully decrypted to: %s\n", outPath)
		}
		return nil
	},
}

func init() {
	decryptCmd.Flags().StringVarP(&outPath, "outpath", "o", "", "output file path (optional)")
	rootCmd.AddCommand(decryptCmd)
}
