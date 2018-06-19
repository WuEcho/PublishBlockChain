package blc

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
)

//UTXO
type Translation struct {
	//事务哈希
	TransHash  []byte
	//输入事务
	TrsIns     []*TranslationInput
	//输出事务
	TrsOuts    []*TranslationOutput
}

//创建CoinBase类型交易
func NewCoinBaseTransaction(address string) *Translation {

  tsInput := &TranslationInput{[]byte{},10,"genesis data"}

  tsOutput := &TranslationOutput{10,address}

  tsCoinBases := &Translation{[]byte{},[]*TranslationInput{tsInput},[]*TranslationOutput{tsOutput}}

  tsCoinBases.HashTransaction()

  return tsCoinBases
}

//设置哈希值
func (tx *Translation)HashTransaction() {
	var result  bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(result.Bytes())
	tx.TransHash = hash[:]
}