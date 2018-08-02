package blc

import (
	"fmt"
	"os"
	"io/ioutil"
	"log"
	"encoding/gob"
	"crypto/elliptic"
	"bytes"
)

const walletFile = "Wallets.dat"

type Wallets struct {
	WalletsMap map[string]*Wallet
}

// 创建钱包集合
func NewWallets() (*Wallets,error) {

	if _,err := os.Stat(walletFile);os.IsNotExist(err) {
		wallets := &Wallets{}
		wallets.WalletsMap = make(map[string]*Wallet)
		return wallets,err
	}

	fileContent,err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallet Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))

	err = decoder.Decode(&wallet)
	if err != nil {
		log.Panic(err)
	}
	return &wallet,nil
}

// 创建一个新钱包
func (w *Wallets) CreateNewWallet()  {

	wallet := NewWallet()
	fmt.Printf("Address：%s\n",wallet.GetAddress())
	w.WalletsMap[string(wallet.GetAddress())] = wallet
    w.SaveWallets()
}

//将钱包信息写入到文件
func (w *Wallets) SaveWallets()  {
	var content bytes.Buffer
	//注册 为了序列化任何类型
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&w)
	if err != nil {
		log.Panic(err)
	}

	//将序列化的数据写入文件，原文件的数据会被覆盖
	err = ioutil.WriteFile(walletFile,content.Bytes(),0644)
	if err != nil {
		log.Panic(err)
	}
}

