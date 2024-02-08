package common

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"strings"
)

func RSASign(data string, priKey string) (string, error) {
	hashMd5 := md5.Sum([]byte(data))
	hashed := hashMd5[:]

	// rsa
	block, _ := pem.Decode([]byte(priKey))
	if block == nil {
		return "", errors.New("private key error")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.MD5, hashed)
	//fmt.Println("消息的签名信息： ", base64.StdEncoding.EncodeToString(signature))
	if err != nil {
		return "", err
	}
	ciphertext := base64.StdEncoding.EncodeToString(signature)
	return ciphertext, nil
}

func RSAVerifySign(data string, pubKey string, signed string) error {
	hashMd5 := md5.Sum([]byte(data))
	hashed := hashMd5[:]

	block, _ := pem.Decode([]byte(pubKey))
	if block == nil {
		return errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)

	decodeString, err := base64.StdEncoding.DecodeString(signed)
	if err != nil {
		return err
	}
	err = rsa.VerifyPKCS1v15(pub, crypto.MD5, hashed, decodeString)
	if err != nil {
		return err
	}

	return nil
}

func RSAEncryptFromPubkey(data []byte, pubkeyBytes []byte) ([]byte, error) {
	pub, err := parsePublicKey(pubkeyBytes)
	if err != nil {
		return nil, err
	}
	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, data, []byte("HBC_MPC"))
	if err != nil {
		return nil, err
	}
	return encryptedBytes, nil
}

func RSADecryptFromPrivkey(encryptedBytes []byte, privkeyBytes []byte) ([]byte, error) {
	priv, err := parsePrivateKey(privkeyBytes)
	if err != nil {
		return nil, err
	}
	plainBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, encryptedBytes, []byte("HBC_MPC"))
	if err != nil {
		return nil, err
	}
	return plainBytes, nil
}

func parsePublicKey(pemData []byte) (*rsa.PublicKey, error) {
	/*pemData, err := hexToDER(hexStr)
	if err != nil {
		fmt.Println("Failed to convert hex to PEM:", err)
		return nil, err
	}*/
	p, _ := pem.Decode(pemData)
	// The key is expected to be in ASN.1 DER format.
	pub, err := x509.ParsePKIXPublicKey(p.Bytes)
	if err != nil {
		pub, err = x509.ParsePKCS1PublicKey(p.Bytes)
		if err != nil {
			return nil, err
		}
	}
	rsaPub := pub.(*rsa.PublicKey)
	return rsaPub, nil
}

// pemDataBytes hex and pem private key
func parsePrivateKey(pemData []byte) (*rsa.PrivateKey, error) {
	// Convert the hex private key to PEM format
	pemDataBytes, err := hexToDER(strings.TrimSpace(string(pemData)))
	if err == nil {
		pemData = pemDataBytes
	}

	// The key is expected to be in ASN.1 DER format.
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("Failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		priv1, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		priv = priv1.(*rsa.PrivateKey)
	}

	return priv, nil
}

// Convert hex-encoded data to DER format (binary)
func hexToDER(hexData string) ([]byte, error) {
	// Decode the hex data
	decoded, err := hex.DecodeString(hexData)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func GenerateRSAKeyPair() error {
	// Generate a new RSA private key with 4096 bits
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		fmt.Println("Error generating RSA private key:", err)
		return err
	}

	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	if err != nil {
		fmt.Println("Error MarshalPKCS8PrivateKey:", err)
		return err
	}
	// Encode the private key to the PEM format
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	}
	privateKeyFile, err := os.Create("./private_key.pem") //"./private_key.pem"
	if err != nil {
		fmt.Println("Error creating private key file:", err)
		return err
	}
	pem.Encode(privateKeyFile, privateKeyPEM)
	privateKeyFile.Close()

	pubASN1 := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	if err != nil {
		fmt.Println("Error MarshalPKIXPublicKey:", err)
		return err
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	publicKeyFile, err := os.Create("./public_key.pem")
	if err != nil {
		fmt.Println("Error creating public key file:", err)
		return err
	}
	defer publicKeyFile.Close()

	_, err = publicKeyFile.Write(pubBytes)
	if err != nil {
		fmt.Println("Error writing public key file:", err)
		return err
	}

	fmt.Println("RSA key pair generated successfully!")
	return nil
}

type Record struct {
	Chain   string
	Path    string
	Address string
}
