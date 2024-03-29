# HBC Recovery Tool
This tool should run in an offline environment

This document is applicable to CoinCover backup

## Getting started


1. install git, go
2. download the project
```
   git clone git@github.com:hblzong/recovery-tool.git
```
3. Compile
```
cd recovery-tool
go mod vendor
go build  -o recovery-tool  main.go
```
4. Get user passphrase, hbc passphrase from HyperBC, and the RSA private key from CoinCover, then run 

    userPassphrase: The customer should remember his own recovery passphrase which he set in the App.

    hbcPassphrase: After customer set the recovery passphrase, the App sends the encrypted key share to Hbc, then Hbc will send him
                  his Hbc recovery passphrase by email. 

    privkeyFilePath: CoinCover would give a RSA private key, copy the whole string into a file privkeyFile, save the file.


The command description:

a. parseZipFile
```
recovery-tool parseZipFile [zipFilePath] [userPassphrase] [hbcPassphrase] [privkeyFilePath]
eg: recovery-tool parseZipFile './zipTest.zip' '123123' '456456' './privkeyFile'
```

Change the ./zipTest.zip to the name of the real backup file that CoinCover provided to you.

The './zipTest.zip' file include the following files:

1. chaincodes_hbc
2. hbc_share.0_hbc
3. hbc_share.1_hbc
4. pubkeys_hbc
5. user_share

Change 123123 to the recovery phrase set in the App.

Change 456456 to the Hbc password in the email.

Change './privkeyFile' to your own RSA private key file.

This should print something like this:

```
{
    "chaincodes": "[\"123123\",\"234234\",\"456456\"]",
    "hbc_share.0": "223344",
    "hbc_share.1": "445566",
    "pubkeys": "[\"123123ab\",\"234234ab\",\"456456ab\"]",
    "user_share": "112233"
}
```

This outputs three arrays:

hbc.encrypted.0/1 and user are private key slices. 

pubkeys is an array of the public key slices.

chaincodes is also an array.

After you get the metadata plaintext, you can use the following three ways to recover the extended child private key(s).

b.deriveChildPrivateKey

```
recovery-tool deriveChildPrivateKey [metadataFilePath] [derivePath]
eg: recovery-tool deriveChildPrivateKey './metadata.json' '81/0/1/60/2'
```
the derivePath is used for derive the child private key, the path is like '81/0/46/0/0'

c.deriveCsvFile

The csvFilePath is From Hbc in the email.
```
recovery-tool deriveCsvFile [metadataFilePath] [csvFilePath]
This command is used for batch recovering. You can download the csv file from hbc backend.
```

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

Metadata: Copy the content in metadata json or the output of parseZipFile, then paste it into this text box.

WalletType: Fund Wallet or Api Wallet.

VaultIndex: If it is Fund Wallet, it starts from 0. Else if it is Api Wallet, this is fixed 0.

Chain: Choose the destination chain from the list.

SubIndex: If it is Fund Wallet, this is fixed 0. Else if it is Api Wallet, it starts from 0.


