package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	var cmdScan = &cobra.Command{
		Use:   "scan [image file]",
		Short: "Scan a QR code image",
		Long: `Scan a QR code image and print its contents to the stdout.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
			fmt.Println("Print: " + args[0])
		},
	}

	var rootCmd = &cobra.Command{Use: "totp"}
	rootCmd.AddCommand(cmdScan)
	rootCmd.Execute()
}
