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
	"github.com/xlzd/gotp"
)

const serviceName = "macOS TOTP CLI"

func addItem(name, secret string) error {
	// Store it to the keychain
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(serviceName)
	item.SetAccount(name)
	item.SetLabel(name)
	item.SetData([]byte(secret))
	item.SetSynchronizable(keychain.SynchronizableNo)
	item.SetAccessible(keychain.AccessibleWhenPasscodeSetThisDeviceOnly)
	return keychain.AddItem(item)
}

func main() {
	var useBarcodeHintWhenScan bool

	var cmdScan = &cobra.Command{
		Use:   "scan <name> <image>",
		Short: "Scan a QR code image",
		Long:  `Scan a QR code image and store it to the macOS keychain.`,
		Args:  cobra.ExactArgs(2),

		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			path := args[1]

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

			var hint map[gozxing.DecodeHintType]interface{}
			if useBarcodeHintWhenScan {
				hint = map[gozxing.DecodeHintType]interface{}{
					gozxing.DecodeHintType_PURE_BARCODE: struct{}{},
				}
			}

			result, err := qrReader.Decode(bmp, hint)
			if err != nil {
				return err
			}

			// parse TOTP URL
			parsed, err := url.Parse(result.GetText())
			if err != nil {
				return err
			}
			secret := parsed.Query().Get("secret")
			// Reference: https://github.com/google/google-authenticator/wiki/Key-Uri-Format
			if parsed.Scheme != "otpauth" || parsed.Host != "totp" || secret == "" {
				return errors.New("Given QR code is not for TOTP")
			}

			// Save to the keychain
			err = addItem(name, secret)
			if err != nil {
				return err
			}
			fmt.Printf("Given QR code successfully registered as \"%v\".\n", name)
			return nil
		},
	}

	cmdScan.Flags().BoolVarP(
		&useBarcodeHintWhenScan,
		"barcode",
		"b",
		false,
		"use PURE_BARCODE hint for decoding. this flag maybe solves FormatException",
	)

	var cmdAdd = &cobra.Command{
		Use:   "add <name>",
		Short: "Manually add a secret to the macOS keychain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			// Read secret from stdin
			var secret string
			fmt.Print("Type secret: ")
			fmt.Scanln(&secret)
			if secret == "" {
				return errors.New("No secret was given")
			}

			// Save to the keychain
			err := addItem(name, secret)
			if err != nil {
				return err
			}
			fmt.Printf("Given secret successfully registered as \"%v\".\n", name)
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

	var cmdGet = &cobra.Command{
		Use:   "get <name>",
		Short: "Get a TOTP code",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			// Query an item
			query := keychain.NewItem()
			query.SetSecClass(keychain.SecClassGenericPassword)
			query.SetService(serviceName)
			query.SetAccount(name)
			query.SetMatchLimit(keychain.MatchLimitOne)
			query.SetReturnData(true)
			results, err := keychain.QueryItem(query)
			if err != nil {
				return err
			}
			if len(results) != 1 {
				return errors.New("Given name is not found")
			}
			r := results[0]

			// Generate a TOTP code
			fmt.Println(gotp.NewDefaultTOTP(string(r.Data)).Now())
			return nil
		},
	}

	var cmdDelete = &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a TOTP code",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			// Delete an item
			query := keychain.NewItem()
			query.SetSecClass(keychain.SecClassGenericPassword)
			query.SetService(serviceName)
			query.SetAccount(name)
			query.SetMatchLimit(keychain.MatchLimitOne)
			err := keychain.DeleteItem(query)
			if err != nil {
				return err
			}

			fmt.Printf("Successfully deleted \"%v\".\n", name)
			return nil
		},
	}

	var rootCmd = &cobra.Command{Use: os.Args[0], Version: "1.0.0"}
	rootCmd.AddCommand(cmdScan, cmdAdd, cmdList, cmdGet, cmdDelete)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
