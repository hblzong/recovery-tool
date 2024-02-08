package common

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"os"
	"strings"
)

// 生成ECC椭圆曲线密钥对，保存到文件
func GenerateECCKey() {
	//生成密钥对
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		panic(err)
	}
	//保存私钥
	//生成文件
	privatefile, err := os.Create("eccprivate.pem")
	if err != nil {
		panic(err)
	}
	//x509编码
	eccPrivateKey, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}
	//pem编码
	privateBlock := pem.Block{
		Type:  "ecc private key",
		Bytes: eccPrivateKey,
	}
	pem.Encode(privatefile, &privateBlock)
	//保存公钥
	publicKey := privateKey.PublicKey
	//创建文件
	publicfile, err := os.Create("eccpublic.pem")
	//x509编码
	eccPublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}
	//pem编码
	block := pem.Block{Type: "ecc public key", Bytes: eccPublicKey}
	pem.Encode(publicfile, &block)
}

// 取得ECC私钥
func GetECCPrivateKey(path string) *ecdsa.PrivateKey {
	//读取私钥
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	//pem解码
	block, _ := pem.Decode(buf)
	//x509解码
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	return privateKey
}

// 取得ECC公钥
func GetECCPublicKey(path string) *ecdsa.PublicKey {
	//读取公钥
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	//pem解密
	block, _ := pem.Decode(buf)
	//x509解密
	publicInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	publicKey := publicInterface.(*ecdsa.PublicKey)
	return publicKey
}

// 对消息的散列值生成数字签名
func SignECC(msg []byte, privateKey *ecdsa.PrivateKey) ([]byte, []byte) {
	//计算哈希值
	hash := sha256.New()
	//填入数据
	hash.Write(msg)
	bytes := hash.Sum(nil)
	//对哈希值生成数字签名
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, bytes)
	if err != nil {
		panic(err)
	}
	rtext, _ := r.MarshalText()
	stext, _ := s.MarshalText()
	return rtext, stext
}

// 验证数字签名
func VerifySignECC(msg []byte, rtext, stext []byte, path string) bool {
	//读取公钥
	publicKey := GetECCPublicKey(path)
	//计算哈希值
	hash := sha256.New()
	hash.Write(msg)
	bytes := hash.Sum(nil)
	//验证数字签名
	var r, s big.Int
	r.UnmarshalText(rtext)
	s.UnmarshalText(stext)
	verify := ecdsa.Verify(publicKey, bytes, &r, &s)
	return verify
}

// 二次封装签名结果为一个字符串结构 "数字_数字"，对应为"r_s"值
func CreateSign(data interface{}, privKey *ecdsa.PrivateKey) (str string) {
	msg, _ := json.Marshal(data)
	rtext, stext := SignECC(msg, privKey)
	str = string(rtext)
	str += "_" + string(stext)
	return str
}

// 解析字符串为R,S两个字节数组
func ParseSignStr(str string) (r, s []byte) {
	if len(str) < 20 {
		return
	}
	pos := strings.Index(str, "_")

	r = []byte(str[:pos])
	s = []byte(str[pos+1:])
	return
}
