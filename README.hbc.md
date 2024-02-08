# HBC Recovery Tool
This tool should run in an offline environment

This document is applicable to HBC backup
## Getting started


1. install git, go
2. download the project
```
   git clone git@james.dumping.ren:mpc/recovery-tool.git
```
3. Compile
```
cd recovery-tool
go mod vendor
go build  -o recovery-tool  main.go
```
4. Get user passphrase and user RSA private key, then run 

    userPassphrase: The customer should remember his own recovery passphrase which he set in the App.

    privkeyFilePath: The customer should remember his own RSA private key, copy the whole string into a file privkeyFile, save the file.


The command description:

a. parseHbcZipFile

Parse the recovery package that HyperBC provided to you.

The command need four inputs:

1. the file path to the recovery package that HyperBC provided to you
2. the file path to you private RSA key 
3. user passphrase that you set in the App

The command should create a file 'metadata.json' and print something like this:

```
{
    "chaincodes": "[\"123123\",\"234234\",\"456456\"]",
    "hbc_share.0": "223344",
    "hbc_share.1": "445566",
    "pubkeys": "[\"123123ab\",\"234234ab\",\"456456ab\"]",
    "user_share": "112233"
}
```

Description:

hbc_share.0/1 and user_share are private key slices.

pubkeys is an array of the public key slices.

chaincodes is also an array.

After you get the metadata plaintext, you can use the following three ways to recover the extended child private key(s).

b. deriveChildPrivateKey

Use metadata plaintext recover the extended child private key by one derivePath.

The command need two inputs:

1. the file path to the metadata.json
2. the derivePath is used for derive the child private key, the path is like '81/0/46/0/0'

c. deriveCsvFile

Use metadata plaintext recover the extended child private keys by a derive path csv file.

The command need two inputs:

1. the file path to the metadata.json
2. the file path to the derive path csv file


d.Use the UI
```
Follow the four steps to generate the derive.html, it is a UI for extended child private key.
1. echo $GOROOT    if it is empty,  execute    export GOROOT= where your go install dir
   cp $GOROOT/misc/wasm/wasm_exec.js .
2. GOOS=js GOARCH=wasm go build -tags=osusergo -o recovery-tool.wasm helpers/helper.go
3. If you are Ubuntu/Linux:  base64 -w 0 recovery-tool.wasm  >  wasmstr.txt
   If you are Mac: base64 -b -i recovery-tool.wasm -o wasmstr.txt
4. awk 'NR==FNR{a[i++]=$0;next} /var base64String = ".*";/{sub(/var base64String = ".*";/, "var base64String = \""a[0]"\";")}1' wasmstr.txt ui/derive_template.html > tmp && mv tmp ui/derive.html
You should install awk before execute this step. 
```

Explain the four steps:

step1 is to move the necessary depencency lib to the project folder

step2 generates recovery-tool.wasm

step3 generates the base64 string of recovery-tool.wasm

step4 replaces the content in template html to the real base64 string, and generates a derive.html file.

Now you can double click the ui/derive.html.

The text boxes or drop down box in the UI:

Metadata: Copy the content in 'metadata.json' or the output of parseZipFile, then paste it into this text box.

WalletType: Fund Wallet or Api Wallet.

VaultIndex: If it is Fund Wallet, it starts from 0. Else if it is Api Wallet, this is fixed 0.

Chain: Choose the destination chain from the list.

SubIndex: If it is Fund Wallet, this is fixed 0. Else if it is Api Wallet, it starts from 0.


