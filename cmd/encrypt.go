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

		fmt.Printf("Successfully encrypted: %s -> %s.genc\n", file, file)

		return nil
	},
}

func init() {
	encryptCmd.MarkPersistentFlagRequired("file")
	encryptCmd.MarkPersistentFlagRequired("passphrase")
	encryptCmd.Flags().BoolVar(&deleteOrigin, "delete-origin", false, "remove original file after encryption")
	rootCmd.AddCommand(encryptCmd)

}
