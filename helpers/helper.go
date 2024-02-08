package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hblzong/recovery-tool/common"
	"strings"
	"syscall/js"
)

func main() {
	js.Global().Set("generateChildExtendedPrivateKey", js.FuncOf(generateChildExtendedPrivateKey))
	js.Global().Set("generateChildExtendedPrivateKeyFromFile", js.FuncOf(generateChildExtendedPrivateKeyFromFile))
	done := make(chan struct{}, 0)
	<-done
}

func generateChildExtendedPrivateKeyFromFile(this js.Value, args []js.Value) interface{} {
	metadataFileStr := args[0].String() //metadata.json文件内容
	walletType := args[1].Float()       //钱包类型 0 资管钱包 1 api钱包
	fileContent := args[2].String()
	paths := strings.Split(fileContent, "\n")

	var metadataMap map[string]string
	err := json.Unmarshal([]byte(metadataFileStr), &metadataMap)
	if err != nil {
		panic(err)
	}

	//参数验证
	walletTypeInt := int(walletType)
	if walletTypeInt != 0 && walletTypeInt != 1 {
		return "wallet type must be 0 or 1"
	}

	result := make([]map[string]interface{}, 0)
	for _, item := range paths {
		itemInfo := strings.Split(item, ",")
		if strings.ToLower(itemInfo[0]) == "chain" {
			continue
		}

		if len(itemInfo) != 3 {
			return fmt.Sprintf("path: %s is invalid", item)
		}

		hdPath := itemInfo[2]

		privBytes, addr, err := common.DeriveChildPrivateKey(metadataMap, hdPath)
		if err != nil {
			panic(err)
		}

		if addr != itemInfo[1] {
			panic(fmt.Errorf("child Extended address: %s not equal file input address: %s", addr, itemInfo[1]))
		}
		result = append(result, map[string]interface{}{
			"private_key": hex.EncodeToString(privBytes),
			"address":     addr,
		})
	}
	resultBytes, _ := json.Marshal(result)
	return string(resultBytes)
}

func generateChildExtendedPrivateKey(this js.Value, args []js.Value) interface{} {
	metadataFileStr := args[0].String() //metadata.json文件内容
	walletType := args[1].Float()       //钱包类型 0 资管钱包 1 api钱包
	vaultIndex := args[2].Float()       //钱包序号，资管钱包从0开始，api钱包固定为0
	chainInt := args[3].Float()         //币种编号
	subIndex := args[4].Float()         //币种下的地址序号，资管钱包固定为0，api钱包从0开始

	//参数验证
	walletTypeInt := int(walletType)
	vaultIndexInt := int(vaultIndex)
	subIndexInt := int(subIndex)
	if walletTypeInt != 0 && walletTypeInt != 1 {
		return "wallet type must be 0 or 1"
	}

	//api钱包 vaultIndex 固定为0
	if walletTypeInt == 1 {
		if vaultIndexInt != 0 {
			return "api wallet vault index must be 1"
		}
	}
	//资产钱包 subIndex 固定为0
	if walletTypeInt == 0 {
		if subIndexInt != 0 {
			return "asset wallet sub index must be 0"
		}
	}

	fmt.Printf("params: %+v", args)
	// 实现生成子扩展私钥的逻辑
	hdPath := fmt.Sprintf("81/%d/%d/%d/%d", walletTypeInt, vaultIndexInt, int(chainInt), subIndexInt)

	fmt.Printf("hdPath: %s", hdPath)
	var metadataMap map[string]string
	err := json.Unmarshal([]byte(metadataFileStr), &metadataMap)
	if err != nil {
		panic(err)
	}
	privBytes, addr, err := common.DeriveChildPrivateKey(metadataMap, hdPath)
	if err != nil {
		panic(err)
	}
	result := make([]map[string]interface{}, 1)
	result[0] = map[string]interface{}{
		"private_key": hex.EncodeToString(privBytes),
		"address":     addr,
	}
	resultBytes, _ := json.Marshal(result)
	return string(resultBytes)
}
