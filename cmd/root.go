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

	"github.com/spf13/cobra"
)

const MinPassLen = 10

var (
	passphrase string
	file       string
)

var rootCmd = &cobra.Command{
	Use:     "go-enc",
	Short:   "Encrypt and decrypt files",
	Long:    `go-enc - A CLI tool for encrypting and decrypting files using AES-256`,
	Version: "1.0",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "path to file (required)")
	rootCmd.PersistentFlags().StringVarP(&passphrase, "passphrase", "p", "", "encryption passphrase (required)")

	// rootCmd.MarkPersistentFlagRequired("file")
	// rootCmd.MarkPersistentFlagRequired("passphrase")
}
