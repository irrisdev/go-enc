/*
Copyright © 2026 irrisdev lithium8260@proton.me

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

		// Validate file has .genc extension
		if !strings.HasSuffix(file, ".genc") {
			return fmt.Errorf("file must have .genc extension: %s", file)
		}

		// Check if file is readable
		f, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("file is not readable: %w", err)
		}
		f.Close()

		// Validate output path if provided
		if outPath != "" {
			// Check if output directory exists
			outDir := filepath.Dir(outPath)
			if outDir != "." && outDir != "" {
				if info, err := os.Stat(outDir); err != nil {
					if os.IsNotExist(err) {
						return fmt.Errorf("output directory does not exist: %s", outDir)
					}
					return fmt.Errorf("error accessing output directory: %w", err)
				} else if !info.IsDir() {
					return fmt.Errorf("output directory path is not a directory: %s", outDir)
				}
			}
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var outputFile string

		if outPath == "" {
			outputFile = strings.TrimSuffix(file, ".genc")
			err = genc.Decrypt(passphrase, file)
		} else {
			outputFile = outPath
			err = genc.Decrypt(passphrase, file, outPath)
		}

		if err != nil {
			return fmt.Errorf("decryption failed: %w", err)
		}

		fmt.Printf("successfully decrypted: %s -> %s\n", file, outputFile)
		return nil
	},
}

func init() {
	decryptCmd.MarkPersistentFlagRequired("file")
	decryptCmd.MarkPersistentFlagRequired("passphrase")
	decryptCmd.Flags().StringVarP(&outPath, "outpath", "o", "", "output file path (optional)")
	rootCmd.AddCommand(decryptCmd)
}
