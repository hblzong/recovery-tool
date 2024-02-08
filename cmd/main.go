package main

import (
	"encoding/hex"
	"fmt"
	"github.com/hblzong/recovery-tool/common"
	"os"
	"strings"
)

//export hello
func hello() string {
	return "Greetings"
}

//export generateChildExtendedPrivateKey
func generateChildExtendedPrivateKey(metadataFilePath string, walletType int, vaultIndex int, chainInt int, subIndex int) string {
	// 实现生成子扩展私钥的逻辑
	hdPath := fmt.Sprintf("81/%d/%d/%d/%d", walletType, vaultIndex, chainInt, subIndex)
	metadataMap, err := common.ReadMetadataFile(metadataFilePath)
	if err != nil {
		panic(err)
	}
	privBytes, _, err := common.DeriveChildPrivateKey(metadataMap, hdPath)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(privBytes)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	if os.Args[1] == "makeZipFile" {
		if len(os.Args) < 10 {
			printUsage()
			os.Exit(1)
		}
		//"2. recovery-tool makeZipFile [zipFilePath] [userPassphrase] [hbcPassphrase] [pubkeyFilePath] [userPrivatekeySlice] [hbcPrivatekeySlice1|hbcPrivatekeySlice2] [chaincode1|chaincode2|chaincode3] [pubkeySlice1|pubkeySlice2|pubkeySlice3]")
		//arg[2]: zipFilePath
		//arg[3]: userPassphrase
		//arg[4]: hbcPassphrase
		//arg[5]: pubkeyFilePath
		//arg[6]: userPrivatekeySlice
		//arg[7]: hbcPrivatekeySlice1|hbcPrivatekeySlice2
		//arg[8]: chaincode1|chaincode2|chaincode3
		//arg[9]: pubkeySlice1|pubkeySlice2|pubkeySlice3
		pubkeyBytes, err := os.ReadFile(os.Args[5])
		if err != nil {
			panic(err)
		}
		//pubkeyBytes = pubkeyBytes[:len(pubkeyBytes)-1]
		/*fmt.Println(string(pubkeyBytes))
		fmt.Println(len(string(pubkeyBytes)))
		fmt.Println(len(pubkeyBytes))
		fmt.Printf("%x \n", pubkeyBytes)*/
		hbcPasswdBytes, err := hex.DecodeString(os.Args[4])
		if err != nil {
			panic(err)
		}
		_, err = common.MakeZipFile([]byte(os.Args[3]), hbcPasswdBytes, pubkeyBytes, os.Args[6], strings.Split(os.Args[7], "|"), strings.Split(os.Args[8], "|"), strings.Split(os.Args[9], "|"), os.Args[2])
		if err != nil {
			panic(err)
		}
		fmt.Println("the zip file is successfully created")
	} else if os.Args[1] == "parseZipFile" {
		if len(os.Args) < 6 {
			printUsage()
			os.Exit(1)
		}
		privkeyBytes, err := os.ReadFile(os.Args[5])
		if err != nil {
			panic(err)
		}

		hbcPasswdBytes, err := hex.DecodeString(os.Args[4])
		if err != nil {
			panic(err)
		}
		d, err := common.ParseFile(os.Args[2], privkeyBytes, []byte(os.Args[3]), hbcPasswdBytes)
		if err != nil {
			panic(err)
		}
		fmt.Println(d)
	} else if os.Args[1] == "parseHbcZipFile" {
		if len(os.Args) < 5 {
			printUsage()
			os.Exit(1)
		}
		privkeyBytes, err := os.ReadFile(os.Args[4])
		if err != nil {
			panic(err)
		}

		d, err := common.ParseFile(os.Args[2], privkeyBytes, []byte(os.Args[3]), nil)
		if err != nil {
			panic(err)
		}
		fmt.Println(d)
	} else if os.Args[1] == "deriveChildPrivateKey" {
		if len(os.Args) < 4 {
			printUsage()
			os.Exit(1)
		}
		metadataMap, err := common.ReadMetadataFile(os.Args[2])
		if err != nil {
			panic(err)
		}
		var addr string
		_, addr, err = common.DeriveChildPrivateKey(metadataMap, os.Args[3])
		if err != nil {
			panic(err)
		}
		fmt.Printf("derived addr: %s \n", addr)
	} else if os.Args[1] == "generateKeyPair" {
		err := common.GenerateRSAKeyPair()
		if err != nil {
			panic(err)
		}
	} else if os.Args[1] == "deriveCsvFile" {
		if len(os.Args) < 4 {
			printUsage()
			os.Exit(1)
		}
		records, err := common.ParseCsv(os.Args[3])
		if err != nil {
			panic(err)
		}
		metadataMap, err := common.ReadMetadataFile(os.Args[2])
		if err != nil {
			panic(err)
		}
		for _, r := range records {
			_, addr, err := common.DeriveChildPrivateKey(metadataMap, r["Path"])
			if err != nil {
				panic(err)
			}
			fmt.Printf("derived addr: %s \n", addr)
		}
	} else {
		printUsage()
	}
	os.Exit(0)
}

func printUsage() {
	fmt.Println("USAGE:")
	fmt.Println("1. recovery-tool generateKeyPair")
	fmt.Println("description: will generate two files: ./private_key.pem and ./public_key.pem")

	fmt.Println("2. recovery-tool makeZipFile [zipFilePath] [userPassphrase] [hbcPassphrase] [pubkeyFilePath] [userPrivatekeySlice] [hbcPrivatekeySlice1|hbcPrivatekeySlice2] [chaincode1|chaincode2|chaincode3] [pubkeySlice1|pubkeySlice2|pubkeySlice3]")
	fmt.Println("eg: recovery-tool makeZipFile './zipTest.zip' '123123' '456456' './public_key.pem' '5ecd00a8164031b61c7c61578137b83d5c0b57d6dbd8617ece480ec9078442c9' '7ecd00a8164031b61c7c61578137b83d5c0b57d6dbd8617ece480ec9078442c9|8ecd00a8164031b61c7c61578137b83d5c0b57d6dbd8617ece480ec9078442c9' '4ecd00a8164031b61c7c61578137b83d5c0b57d6dbd8617ece480ec9078442c9|4ecd00a8164031b61c7c61578137b83d5c0b57d6dbd8617ece480ec9078442ca|4ecd00a8164031b61c7c61578137b83d5c0b57d6dbd8617ece480ec9078442cb' '033669d206967b384d588b366b6400733987befc6604fec764f9fc2d42a3bf7a86|021b491468a9c042e6d4e994c3979df14454cd99e4fc207161a929e719f644540b|02f75cebd23a9ac7e1364d0462be378f09aaf26474eb46cc43bdef5de2817932e5'")

	fmt.Println("3. recovery-tool parseZipFile [zipFilePath] [userPassphrase] [hbcPassphrase] [privkeyFilePath]")
	fmt.Println("eg: recovery-tool parseZipFile './zipTest.zip' '123123' '456456' './private_key.pem'")

	fmt.Println("4. recovery-tool deriveChildPrivateKey [metadataFilePath] [derivePath]")
	fmt.Println("eg: recovery-tool deriveChildPrivateKey './metadata.json' '81/0/1/60/2'")

	fmt.Println("5. recovery-tool deriveCsvFile [metadataFilePath] [csvFilePath]")
	fmt.Println("eg: recovery-tool deriveCsvFile './metadata.json' './walletPaths.csv'")

	fmt.Println("6. recovery-tool parseHbcZipFile [zipFilePath] [userPassphrase] [privkeyFilePath]")
	fmt.Println("eg: recovery-tool parseHbcZipFile './zipTest.zip' '123123' './private_key.pem'")
}
