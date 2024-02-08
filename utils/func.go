package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/scrypt"
)

func FormatTaskData(task map[string]interface{}) map[string]interface{} {
	start_time := int64(0)
	if _, ok := task["start_time"].(float64); ok {
		start_time = int64(task["start_time"].(float64))
	} else {
		start_time = task["start_time"].(int64)
	}
	send_id := ""
	switch item := task["send_id"].(type) {
	case string:
		_send_id, _ := decimal.NewFromString(item)
		__send_id, _ := _send_id.Float64()
		send_id = decimal.NewFromFloat(__send_id).String()
	case int64:
		send_id = fmt.Sprintf("%d", item)
	case float64:
		send_id = decimal.NewFromFloat(item).String()
	}

	shop_id := ""
	switch item := task["shop_id"].(type) {
	case string:
		_shop_id, _ := decimal.NewFromString(item)
		__shop_id, _ := _shop_id.Float64()
		shop_id = decimal.NewFromFloat(__shop_id).String()
	case int64:
		shop_id = fmt.Sprintf("%d", item)
	case float64:
		shop_id = decimal.NewFromFloat(item).String()
	}

	amount := ""
	switch item := task["amount"].(type) {
	case string:
		_amount, _ := decimal.NewFromString(item)
		__amount, _ := _amount.Float64()
		amount = decimal.NewFromFloat(__amount).String()
	case float64:
		amount = decimal.NewFromFloat(item).String()
	case int64:
		amount = fmt.Sprintf("%d", item)
	}

	teamId := int64(0)
	if _, ok := task["team_id"].(float64); ok {
		teamId = int64(task["team_id"].(float64))
	} else {
		teamId = task["team_id"].(int64)
	}

	walletId := int64(0)
	if _, ok := task["wallet_id"].(float64); ok {
		walletId = int64(task["wallet_id"].(float64))
	} else {
		walletId = task["wallet_id"].(int64)
	}

	walletType := int(0)
	if _, ok := task["wallet_type"].(float64); ok {
		walletType = int(task["wallet_type"].(float64))
	} else {
		walletType = task["wallet_type"].(int)
	}

	data := map[string]interface{}{
		"shop_id":     shop_id,
		"coin":        task["coin"].(string),
		"contract":    task["contract"].(string),
		"start_time":  start_time,
		"send_id":     send_id,
		"team_id":     teamId,
		"wallet_id":   walletId,
		"wallet_type": walletType,
		"amount":      amount,
		"memo":        task["memo"].(string),
		"salt":        task["salt"].(string),
	}
	return data
}

func Aes128Encrypt(key, iv string, src string) (string, error) {
	_key, _iv, err := parseByte(key, iv, 16)
	if err != nil {
		return "", err
	}
	return AesEncrypt(_key, _iv, src), nil
}

func Aes128Decrypt(key, iv string, src string) (string, error) {
	_key, _iv, err := parseByte(key, iv, 16)
	if err != nil {
		return "", err
	}
	return AesDecrypt(_key, _iv, src), nil
}

func Aes256Encrypt(key, iv string, src string) (string, error) {
	_key, _iv, err := parseByte(key, iv, 32)
	if err != nil {
		return "", err
	}
	return AesEncrypt(_key, _iv, src), nil
}

func Aes256Decrypt(key, iv string, src string) (string, error) {
	_key, _iv, err := parseByte(key, iv, 32)
	if err != nil {
		return "", err
	}
	return AesDecrypt(_key, _iv, src), nil
}

func parseByte(key, iv string, bitLen int) (_key []byte, _iv []byte, err error) {
	_key = []byte(key)
	_iv = []byte(iv)

	if len(_iv) != 16 {
		return nil, nil, errors.New("iv length is not 16")
	}

	if len(_key) > bitLen {
		_key = _key[:bitLen]
	} else {
		keyByte := make([]byte, bitLen, bitLen)
		//switch bitLen {
		//case 16:
		//	keyByte = [16]byte{}
		//case 32:
		//	keyByte = [32]byte{}
		//default:
		//	return nil, nil, errors.New("bitLen length is  16 or 32")
		//}
		//转换key为字节数组
		keyByteTemp := key
		//依次赋值,这里的3为key的len
		for i := 0; i < len(_key); i++ {
			keyByte[i] = keyByteTemp[i]
		}
		_key = keyByte[:]
	}
	return
}

// AES加密
func AesEncrypt(key []byte, iv []byte, src string) string {

	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}
	blockSize := block.BlockSize()
	_src := PKCS7Padding([]byte(src), blockSize)

	encryptData := make([]byte, len(_src))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encryptData, _src)
	return base64.StdEncoding.EncodeToString(encryptData)
}

// AES解密
func AesDecrypt(key []byte, iv []byte, secret string) string {
	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}
	src, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return ""
	}
	dst := make([]byte, len(src))

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(dst, src)

	dst = PKCS7UnPadding(dst)
	return string(dst)
}

func AesGcmEncrypt(key, data []byte) ([]byte, error) {
	key, salt, err := deriveKey(key, nil)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	ciphertext = append(ciphertext, salt...)

	return ciphertext, nil
}

func AesGcmDecrypt(key, data []byte) ([]byte, error) {
	salt, data := data[len(data)-32:], data[:len(data)-32]

	key, _, err := deriveKey(key, salt)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func deriveKey(password, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}

	key, err := scrypt.Key(password, salt, 128, 2, 1, 32)
	if err != nil {
		return nil, nil, err
	}

	return key, salt, nil
}

func PKCS7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS7UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
