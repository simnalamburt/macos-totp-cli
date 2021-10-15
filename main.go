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

	keychain "github.com/keybase/go-keychain"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/spf13/cobra"
)

const serviceName = "macOS TOTP CLI"

func main() {
	var cmdScan = &cobra.Command{
		Use:   "scan <path of the QR image> <service name>",
		Short: "Scan a QR code image",
		Long:  `Scan a QR code image and store it to the macOS keychain.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Infer service name from the QR codes
			path := args[0]
			name := args[1]

			// open and decode image file
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

			// Store it to the keychain
			item := keychain.NewItem()
			item.SetSecClass(keychain.SecClassGenericPassword)
			item.SetService(serviceName)
			item.SetAccount(name)
			item.SetLabel(name)
			item.SetData([]byte(secret))
			item.SetSynchronizable(keychain.SynchronizableNo)
			item.SetAccessible(keychain.AccessibleWhenPasscodeSetThisDeviceOnly)
			err = keychain.AddItem(item)
			if err != nil {
				return err
			}

			fmt.Printf("Given QR code successfully registered as \"%v\".\n", name)
			return nil
		},
	}

	var cmdList = &cobra.Command{
		Use:   "list",
		Short: "List all registered TOTP codes",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Query items
			query := keychain.NewItem()
			query.SetSecClass(keychain.SecClassGenericPassword)
			query.SetService(serviceName)
			query.SetMatchLimit(keychain.MatchLimitAll)
			query.SetReturnAttributes(true)
			results, err := keychain.QueryItem(query)
			if err != nil {
				return err
			}

			// List query results
			for _, r := range results {
				fmt.Println(r.Account)
			}
			return nil
		},
	}

	var rootCmd = &cobra.Command{Use: os.Args[0]}
	rootCmd.AddCommand(cmdScan, cmdList)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
