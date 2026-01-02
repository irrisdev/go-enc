package cmd

import (
	"github.com/irrisdev/go-enc/genc"
	"github.com/spf13/cobra"
)

var outPath string

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt a file",
	Run: func(cmd *cobra.Command, args []string) {
		if outPath == "" {
			genc.Decrypt(passphrase, file)
		} else {
			genc.Decrypt(passphrase, file, outPath)
		}
	},
}

func init() {
	decryptCmd.Flags().StringVarP(&outPath, "outpath", "o", "", "output file path")
	rootCmd.AddCommand(decryptCmd)
}
