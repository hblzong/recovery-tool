package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hblzong/recovery-tool/common"
	"gopkg.in/urfave/cli.v1"
	"os"
)

func main() {
	app := &cli.App{}
	app.Name = "HyperBC Recovery Tool"
	app.Usage = "A tool for recovery your private key."
	app.UsageText = "recovery-tool [global options] command [arguments...]"
	app.Version = "0.0.1"
	app.BashComplete = cli.DefaultAppComplete

	app.Commands = []cli.Command{
		{
			Name:      "generateRSAKeyPair",
			ShortName: "rsa",
			Usage:     "Will generate two files: ./private_key.pem and ./public_key.pem",
			Action: func(c *cli.Context) error {
				err := common.GenerateRSAKeyPair()
				if err != nil {
					fmt.Printf("generateRSAKeyPair error: %s!\n", err.Error())
				}
				return nil
			},
		},
		{
			Name:      "parseZipFile",
			ShortName: "pzf",
			Usage:     "Parse the recovery package that CoinCover provided to you.",
			UsageText: "",
			Action: func(c *cli.Context) error {
				var zipFilePath string
				fmt.Println("Please input the file path to the recovery package that CoinCover provided to you:")
				if _, err := fmt.Scanln(&zipFilePath); err != nil {
					fmt.Println("input error", err)
					return nil
				}

				if _, err := os.Stat(zipFilePath); os.IsNotExist(err) {
					fmt.Println("no such file")
					return nil
				}

				var privkeyFile string
				fmt.Println("Please input the file path to the private RSA key that CoinCover provided you:")
				if _, err := fmt.Scanln(&privkeyFile); err != nil {
					fmt.Println("input error", err)
					return nil
				}

				privkeyBytes, err := os.ReadFile(privkeyFile)
				if err != nil {
					fmt.Printf("read file: %s error, %s\n", privkeyFile, err.Error())
					return nil
				}

				var userPassphrase string
				fmt.Println("Please input user passphrase which you set in the App:")
				if _, err = fmt.Scanln(&userPassphrase); err != nil {
					fmt.Println("input error", err)
					return nil
				}

				var hbcPassphrase string
				fmt.Println("Please input hbc passphrase than hbc send you by email:")
				if _, err = fmt.Scanln(&hbcPassphrase); err != nil {
					fmt.Println("input error", err)
					return nil
				}
				hbcPasswdBytes, err := hex.DecodeString(hbcPassphrase)
				if err != nil {
					fmt.Println("hbc passphrase hex decode error")
					return nil
				}
				d, err := common.ParseFile(zipFilePath, privkeyBytes, []byte(userPassphrase), hbcPasswdBytes)
				if err != nil {
					fmt.Printf("Parse File error, %s\n", err.Error())
					return nil
				}
				fmt.Println("Parse ZIP File success and metadata.json created")
				dataBytes, _ := json.MarshalIndent(d, "", "    ")
				fmt.Printf("%s\n", dataBytes)

				return nil
			},
		},
		{
			Name:      "deriveChildPrivateKey",
			ShortName: "dcp",
			Usage:     "Use metadata plaintext recover the extended child private key.",
			UsageText: "",
			Action: func(c *cli.Context) error {
				var metadataPath string
				fmt.Println("Please input the metadata file path (default: './metadata.json'):")
				if _, err := fmt.Scanln(&metadataPath); err != nil {
					if len(metadataPath) <= 0 {
						metadataPath = "./metadata.json"
						fmt.Printf("%s\n", metadataPath)
					} else {
						fmt.Printf("input error, %s\n", err.Error())
						return nil
					}
				}

				if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
					fmt.Println("no such file")
					return nil
				}

				var hdPath string
				fmt.Println("Please input a derivePath is used for derive the child private key(eg: 81/0/1/60/2):")
				if _, err := fmt.Scanln(&hdPath); err != nil {
					fmt.Printf("input error, %s\n", err.Error())
					return nil
				}

				metadataMap, err := common.ReadMetadataFile(metadataPath)
				if err != nil {
					fmt.Printf("read metadata file error, %s\n", err.Error())
					return nil
				}

				_, _, err = common.DeriveChildPrivateKey(metadataMap, hdPath)
				if err != nil {
					fmt.Printf("derive child privateKey error, %s\n", err.Error())
					return nil
				}

				return nil
			},
		},
		{
			Name:      "deriveCsvFile",
			ShortName: "dcf",
			Usage:     "Use metadata plaintext recover the extended child private key from csv file.",
			UsageText: "",
			Action: func(c *cli.Context) error {
				var metadataPath string
				fmt.Println("Please input the metadata file path (default: './metadata.json'):")
				if _, err := fmt.Scanln(&metadataPath); err != nil {
					if len(metadataPath) <= 0 {
						metadataPath = "./metadata.json"
						fmt.Printf("%s\n", metadataPath)
					} else {
						fmt.Printf("input error, %s\n", err.Error())
						return nil
					}
				}

				if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
					fmt.Println("no such file")
					return nil
				}

				var csvFilePath string
				fmt.Println("Please input the csv file path:")
				if _, err := fmt.Scanln(&csvFilePath); err != nil {
					fmt.Printf("input error, %s\n", err.Error())
					return nil
				}

				if _, err := os.Stat(csvFilePath); os.IsNotExist(err) {
					fmt.Println("no such file")
					return nil
				}

				metadataMap, err := common.ReadMetadataFile(metadataPath)
				if err != nil {
					fmt.Printf("read metadata file error, %s\n", err.Error())
					return nil
				}

				records, err := common.ParseCsv(csvFilePath)
				if err != nil {
					fmt.Printf("parse csv file error, %s\n", err.Error())
					return nil
				}

				for _, r := range records {
					_, _, err = common.DeriveChildPrivateKey(metadataMap, r["Path"])
					if err != nil {
						fmt.Printf("derive child privateKey error, %s\n", err.Error())
						return nil
					}
				}

				return nil
			},
		},
		{
			Name:      "parseHbcZipFile",
			ShortName: "phzf",
			Usage:     "Parse the recovery package that HyperBc provided to you.",
			UsageText: "",
			Action: func(c *cli.Context) error {
				var zipFilePath string
				fmt.Println("Please input the file path to the recovery package that HyperBc provided to you:")
				if _, err := fmt.Scanln(&zipFilePath); err != nil {
					fmt.Println("input error", err)
					return nil
				}

				if _, err := os.Stat(zipFilePath); os.IsNotExist(err) {
					fmt.Println("no such file")
					return nil
				}

				var privkeyFile string
				fmt.Println("Please input the file path to the private RSA key that you have:")
				if _, err := fmt.Scanln(&privkeyFile); err != nil {
					fmt.Println("input error", err)
					return nil
				}

				privkeyBytes, err := os.ReadFile(privkeyFile)
				if err != nil {
					fmt.Printf("read file: %s error, %s\n", privkeyFile, err.Error())
					return nil
				}

				var userPassphrase string
				fmt.Println("Please input user passphrase which you set in the App:")
				if _, err = fmt.Scanln(&userPassphrase); err != nil {
					fmt.Println("input error", err)
					return nil
				}

				d, err := common.ParseFile(zipFilePath, privkeyBytes, []byte(userPassphrase), nil)
				if err != nil {
					fmt.Printf("Parse File error, %s\n", err.Error())
					return nil
				}
				fmt.Println("Parse HBC ZIP File success and metadata.json created")
				dataBytes, _ := json.MarshalIndent(d, "", "    ")
				fmt.Printf("%s\n", dataBytes)

				return nil
			},
		},
	}

	app.Run(os.Args)
}
