package main

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/url"
	"os"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/spf13/cobra"
)

func main() {
	var cmdScan = &cobra.Command{
		Use:   "scan [image file]",
		Short: "Scan a QR code image",
		Long:  `Scan a QR code image and store it to the macOS keychain.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// open and decode image file
			path := args[0]
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			img, _, err := image.Decode(file)
			if err != nil {
				return err
			}

			// prepare BinaryBitmap
			bmp, err := gozxing.NewBinaryBitmapFromImage(img)
			if err != nil {
				return err
			}

			// decode image
			qrReader := qrcode.NewQRCodeReader()
			result, err := qrReader.Decode(bmp, nil)
			if err != nil {
				return err
			}

			// parse TOTP URL
			parsed, err := url.Parse(result.GetText())
			if err != nil {
				return err
			}
			secret := parsed.Query().Get("secret")
			if parsed.Scheme != "otpauth" || parsed.Host != "totp" || secret == "" {
				return errors.New("Given QR code is not for TOTP")
			}

			// TODO: Store it to the keychain
			fmt.Println(secret)

			return nil
		},
	}

	var rootCmd = &cobra.Command{Use: "totp"}
	rootCmd.AddCommand(cmdScan)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
