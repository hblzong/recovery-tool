/**
*支持aes128、aes256加解密
*参考文档
*https://blog.csdn.net/a_lzq/article/details/108204967?utm_medium=distribute.pc_relevant.none-task-blog-BlogCommendFromBaidu-19.control&dist_request_id=1328626.22722.16154492627822725&depth_1-utm_source=distribute.pc_relevant.none-task-blog-BlogCommendFromBaidu-19.control
*
 */
package common

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/md5"
	"crypto/tls"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/hblzong/recovery-tool/tss/crypto"
	"github.com/hblzong/recovery-tool/tss/crypto/ckd"
	"github.com/hblzong/recovery-tool/utils"
	"io"
	"io/ioutil"
	"math/big"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// mnemonic language
const (
	English            = "english"
	ChineseSimplified  = "chinese_simplified"
	ChineseTraditional = "chinese_traditional"
	Korean             = "korean"
)

// zero is deafult of uint32
const (
	Zero      uint32 = 0
	ZeroQuote uint32 = 0x80000000
	BTCToken  uint32 = 0x10000000
	ETHToken  uint32 = 0x20000000
)

// wallet type from bip44
const (
	// https://github.com/satoshilabs/slips/blob/master/slip-0044.md#registered-coin-types
	BTC       = Zero + 0
	LTC       = Zero + 2
	DOGE      = Zero + 3
	DASH      = Zero + 5
	Optimism  = Zero + 10
	ETH       = Zero + 60
	BCH       = Zero + 145
	TRX       = Zero + 195
	BSV       = Zero + 236
	Fantom    = Zero + 250
	ZKSYNC    = Zero + 324
	POLYGON   = Zero + 966
	ARBITRUM  = Zero + 9001
	OKChain   = Zero + 996
	BSC       = Zero + 714
	HECO      = Zero + 553
	Avalanche = Zero + 43114
	Apt       = Zero + 637
	SUI       = Zero + 784
	SOL       = Zero + 501
)

// 全局 transport
var globalTransport *http.Transport

func init() {
	globalTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

func HttpPost(api string, data map[string]interface{}) (code int, body string, err error) {
	req := make(url.Values)
	for key, item := range data {
		req[key] = []string{fmt.Sprintf("%v", item)}
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	//把post表单发送给目标服务器
	res, err := client.PostForm(api, req)
	if err != nil {
		return 0, "", err
	}
	defer res.Body.Close()
	bytess, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, "", err
	}
	return res.StatusCode, string(bytess), nil
}

func HttpGet(api string) (code int, body string, err error) {
	res, err := http.Get(api)
	if err != nil {
		return 0, "", err
	}
	defer res.Body.Close()
	bytess, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, "", err
	}
	return res.StatusCode, string(bytess), nil
}

func Md5Encode(str string) string {
	data := []byte(str)
	h := md5.New()
	h.Write(data)
	output := h.Sum(nil)
	return fmt.Sprintf("%x", output)
}

func Request(params map[string]interface{}, headers map[string]string, url string, res interface{}) error {
	ctx := context.Background()

	var err error
	req := &http.Request{}
	if len(params) > 0 {
		postData, err := json.Marshal(params)
		if err != nil {
			return err
		}
		req, err = http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(postData))
	} else {
		req, err = http.NewRequestWithContext(ctx, "GET", url, nil)
	}

	if err != nil {
		return err
	}

	// headers
	for key, val := range headers {
		req.Header.Add(key, val)
	}

	// http client and send request
	httpclient := &http.Client{
		Transport: globalTransport,
		Timeout:   10 * time.Second,
	}
	resp, err := httpclient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// parse body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(body) != 0 {
		if err := json.Unmarshal(body, &res); err != nil {
			return err
		}
	}

	// return result
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		return fmt.Errorf("get status code: %d", resp.StatusCode)
	}
	return nil
}

func Upload(url string, headers map[string]string, values map[string]io.Reader) (err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	httpclient := &http.Client{
		Transport: globalTransport,
		Timeout:   10 * time.Second,
	}
	// Submit the request
	res, err := httpclient.Do(req)
	if err != nil {
		return
	}

	// Check the response
	if res.StatusCode != http.StatusOK || res.StatusCode != http.StatusCreated {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	return
}

func MakeZipFile(userPassphrase, hbcPassphrase []byte, pubkeyPem []byte, userPrivateSlice string, hbcPrivateSlices []string, chaincodes []string, ownerPubkeySlices []string, saveFilePath string) (*os.File, error) {
	archive, err := os.Create(saveFilePath)
	if err != nil {
		return nil, err
	}
	defer archive.Close()
	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()
	w1, err := zipWriter.Create("user_share")
	if err != nil {
		return nil, err
	}
	encryptedPrivkeySlice, err := utils.AesGcmEncrypt(userPassphrase, []byte(userPrivateSlice))
	if err != nil {
		return nil, err
	}
	encryptedPrivkeySlice, err = RSAEncryptFromPubkey(encryptedPrivkeySlice, pubkeyPem)
	if err != nil {
		return nil, err
	}
	bs := bytes.NewBufferString(hex.EncodeToString(encryptedPrivkeySlice))
	if _, err = io.Copy(w1, bs); err != nil {
		return nil, err
	}

	fileName := "chaincodes"
	if len(hbcPassphrase) == 0 {
		fileName += "_hbc"
	}
	w1, err = zipWriter.Create(fileName)
	if err != nil {
		return nil, err
	}
	chaincodesBytes, err := json.Marshal(chaincodes)
	if err != nil {
		return nil, err
	}
	if len(hbcPassphrase) > 0 {
		encryptedPrivkeySlice, err = utils.AesGcmEncrypt(hbcPassphrase, chaincodesBytes)
		if err != nil {
			return nil, err
		}
	}

	encryptedPrivkeySlice, err = RSAEncryptFromPubkey(encryptedPrivkeySlice, pubkeyPem)
	if err != nil {
		return nil, err
	}
	bs = bytes.NewBufferString(hex.EncodeToString(encryptedPrivkeySlice))
	if _, err = io.Copy(w1, bs); err != nil {
		return nil, err
	}

	fileName = "pubkeys"
	if len(hbcPassphrase) == 0 {
		fileName += "_hbc"
	}
	w1, err = zipWriter.Create(fileName)
	if err != nil {
		return nil, err
	}
	pubkeysBytes, err := json.Marshal(ownerPubkeySlices)
	if err != nil {
		return nil, err
	}
	if len(hbcPassphrase) > 0 {
		encryptedPrivkeySlice, err = utils.AesGcmEncrypt(hbcPassphrase, pubkeysBytes)
		if err != nil {
			return nil, err
		}
	}
	encryptedPrivkeySlice, err = RSAEncryptFromPubkey(encryptedPrivkeySlice, pubkeyPem)
	if err != nil {
		return nil, err
	}
	bs = bytes.NewBufferString(hex.EncodeToString(encryptedPrivkeySlice))
	if _, err = io.Copy(w1, bs); err != nil {
		return nil, err
	}

	for i := 0; i < len(hbcPrivateSlices); i++ {
		fileName = fmt.Sprintf("hbc_share.%d", i)
		if len(hbcPassphrase) == 0 {
			fileName += "_hbc"
		}
		w1, err = zipWriter.Create(fileName)
		if err != nil {
			return nil, err
		}
		if len(hbcPassphrase) > 0 {
			encryptedPrivkeySlice, err = utils.AesGcmEncrypt(hbcPassphrase, []byte(hbcPrivateSlices[i]))
			if err != nil {
				return nil, err
			}
		}
		encryptedPrivkeySlice, err = RSAEncryptFromPubkey(encryptedPrivkeySlice, pubkeyPem)
		if err != nil {
			return nil, err
		}
		bs = bytes.NewBufferString(hex.EncodeToString(encryptedPrivkeySlice))
		if _, err = io.Copy(w1, bs); err != nil {
			return nil, err
		}
	}
	return archive, nil
}

func readAll(file *zip.File) ([]byte, error) {
	fc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fc.Close()

	content, err := ioutil.ReadAll(fc)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func ParseFile(zipFilePath string, privkeyPem []byte, userPassphrase, hbcPassphrase []byte) (map[string]string, error) {
	zf, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return nil, err
	}
	defer zf.Close()

	dataMap := make(map[string]string)
	for _, file := range zf.File {
		fileBytes, err := readAll(file)
		if err != nil {
			return nil, err
		}
		fileContent := string(fileBytes)
		textBytes, err := hex.DecodeString(fileContent)
		if err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
		//fmt.Printf("=%s\n", file.Name)
		//fmt.Printf("%x\n\n", fileBytes) // file content
		encryptedBytes, err := RSADecryptFromPrivkey(textBytes, privkeyPem)
		if err != nil {
			return nil, err
		}
		if strings.Contains(file.Name, "user_share") {
			plainBytes, err := utils.AesGcmDecrypt(userPassphrase, encryptedBytes)
			if err != nil {
				return nil, err
			}
			dataMap["user_share"] = hex.EncodeToString(plainBytes)
		} else {
			if len(hbcPassphrase) > 0 {
				plainBytes, err := utils.AesGcmDecrypt(hbcPassphrase, encryptedBytes)
				if err != nil {
					return nil, err
				}
				if strings.Contains(file.Name, "share") {
					dataMap[strings.TrimRight(file.Name, "_hbc")] = hex.EncodeToString(plainBytes)
				} else {
					dataMap[strings.TrimRight(file.Name, "_hbc")] = string(plainBytes)
				}
			} else {
				if strings.Contains(file.Name, "share") {
					dataMap[file.Name] = hex.EncodeToString(encryptedBytes)
				} else {
					dataMap[file.Name] = string(encryptedBytes)
				}
			}
		}
	}
	file, err := os.Create("./metadata.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return nil, err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(dataMap)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return nil, err
	}

	return dataMap, nil
}

func derivationChildKey(prByte, pubByte, codeByte []byte, path string) (childPrivKey [32]byte, childPubKey []byte, err error) {
	pubkey, err := btcec.ParsePubKey(pubByte, btcec.S256())
	if err != nil {
		return childPrivKey, nil, fmt.Errorf("derive child pubkey err: %s", err.Error())
	}
	ecPoint, err := crypto.NewECPoint(btcec.S256(), pubkey.X, pubkey.Y)
	if err != nil {
		return childPrivKey, nil, fmt.Errorf("derive child private err: %s", err.Error())
	}

	extendedKey := ckd.NewExtendKey(prByte, ecPoint, ecPoint, 0, 0, codeByte)

	childPrivKey, childPubKey, err = ckd.DerivePrivateKeyForPath(extendedKey, path)
	if err != nil {
		return childPrivKey, nil, fmt.Errorf("derive child private err: %s", err.Error())
	}
	return childPrivKey, childPubKey, nil
}

func ReadMetadataFile(metadataFilePath string) (map[string]string, error) {
	file, err := os.Open(metadataFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var params map[string]string
	err = decoder.Decode(&params)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}
	//fmt.Println("metadata:", params)
	return params, nil
}

func TestDeriveChildPrivateKey() {
	derivationChildKey(nil, nil, nil, "")
}

func DeriveChildPrivateKey(params map[string]string, hdPath string) ([]byte, string, error) { //TODO: derive address according to the chain type
	userPrivKeyStr, ok := params["user_share"]
	if !ok {
		panic("invalid zip file, does not contain user_share")
	}
	hbcShare0Str, ok := params["hbc_share.0"]
	if !ok {
		panic("invalid zip file, does not contain user_share")
	}
	hbcShare1Str, ok := params["hbc_share.1"]
	if !ok {
		panic("invalid zip file, does not contain user_share")
	}
	chainCodeStr, ok := params["chaincodes"]
	if !ok {
		panic("invalid zip file, does not contain chaincodes")
	}
	var chainCodes []string
	err := json.Unmarshal([]byte(chainCodeStr), &chainCodes)
	if err != nil {
		return nil, "", err
	}
	pubkeyStr, ok := params["pubkeys"]
	if !ok {
		panic("invalid zip file, does not contain pubkeys")
	}
	var pubkeys []string
	err = json.Unmarshal([]byte(pubkeyStr), &pubkeys)
	if err != nil {
		return nil, "", err
	}
	privateKeyBytes, err := hex.DecodeString(userPrivKeyStr)
	if err != nil {
		return nil, "", err
	}
	chainCodeBytes, err := hex.DecodeString(chainCodes[0])
	if err != nil {
		return nil, "", err
	}
	deducePubkeyBytes, err := hex.DecodeString(pubkeys[0])
	if err != nil {
		return nil, "", err
	}
	childPrivateKeySlice, _, err := derivationChildKey(privateKeyBytes, deducePubkeyBytes, chainCodeBytes, hdPath)
	if err != nil {
		return nil, "", err
	}

	privateKey := big.NewInt(0).SetBytes(childPrivateKeySlice[:])

	privateKeyBytes, err = hex.DecodeString(hbcShare0Str)
	if err != nil {
		return nil, "", err
	}
	chainCodeBytes, err = hex.DecodeString(chainCodes[1])
	if err != nil {
		return nil, "", err
	}
	deducePubkeyBytes, err = hex.DecodeString(pubkeys[1])
	if err != nil {
		return nil, "", err
	}
	childPrivateKeySlice, _, err = derivationChildKey(privateKeyBytes, deducePubkeyBytes, chainCodeBytes, hdPath)
	if err != nil {
		return nil, "", err
	}
	privateKey.Add(privateKey, big.NewInt(0).SetBytes(childPrivateKeySlice[:]))
	privateKey.Mod(privateKey, btcec.S256().Params().N)

	privateKeyBytes, err = hex.DecodeString(hbcShare1Str)
	if err != nil {
		return nil, "", err
	}
	chainCodeBytes, err = hex.DecodeString(chainCodes[2])
	if err != nil {
		return nil, "", err
	}
	deducePubkeyBytes, err = hex.DecodeString(pubkeys[2])
	if err != nil {
		return nil, "", err
	}
	childPrivateKeySlice, _, err = derivationChildKey(privateKeyBytes, deducePubkeyBytes, chainCodeBytes, hdPath)
	if err != nil {
		return nil, "", err
	}
	privateKey.Add(privateKey, big.NewInt(0).SetBytes(childPrivateKeySlice[:]))
	privateKey.Mod(privateKey, btcec.S256().Params().N)
	pubECPoint := crypto.ScalarBaseMult(btcec.S256(), big.NewInt(0).SetBytes(privateKey.Bytes()))
	publicKey := &ecdsa.PublicKey{X: big.NewInt(0).SetBytes(pubECPoint.X().Bytes()), Y: big.NewInt(0).SetBytes(pubECPoint.Y().Bytes()), Curve: btcec.S256()}

	eles := strings.Split(hdPath, "/")
	chainIntStr := eles[3]
	chainInt, err := strconv.Atoi(chainIntStr)
	if err != nil {
		return nil, "", err
	}
	chainUint32 := uint32(chainInt)
	var chain string
	switch chainUint32 + Zero {
	case BTC:
		chain = "btc"
	case LTC:
		chain = "ltc"
	case DOGE:
		chain = "doge"
	case ETH:
		chain = "eth"
	case BCH:
		chain = "bch"
	//case DASH:
	//	chain = "dash"
	case TRX:
		chain = "trx"
	//eth系列共用地址
	//case HECO:
	//	chain = "eth"
	//case BSC:
	//	chain = "eth"
	default:
		panic("invalid chain type")
	}
	addr, err := switchChainAddress(publicKey, chain)
	if err != nil {
		return nil, "", err
	}
	fmt.Printf("derive child private key success by: %s. \npriv key is: %x \npubkey x: %x y: %x \naddr: %s \n\n", hdPath, privateKey.Bytes(), pubECPoint.X().Bytes(), pubECPoint.Y().Bytes(), addr)
	return privateKey.Bytes(), addr, nil
}

func ParseCsv(filePath string) ([]map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return nil, err
	}

	var data []map[string]string
	for i, row := range records {
		if i == 0 {
			continue //ignore the first row, it is the field header
		}
		record := map[string]string{
			"Chain":   row[0],
			"Address": row[1],
			"Path":    row[2],
		}
		data = append(data, record)
	}

	for _, record := range data {
		fmt.Printf("Chain: %s, Path: %s, Address: %s\n", record["Chain"], record["Path"], record["Address"])
	}

	return data, nil
}
