package blc

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
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


func NewSimpleTranslation(from string,to string,amount int) *Translation {
	var tsIns []*TranslationInput
	var tsOuts  []*TranslationOutput

	bytess,_ := hex.DecodeString("cea12d33b2e7083221bf3401764fb661fd6c34fab50f5460e77628c42ca0e92b")

	tsInput := &TranslationInput{bytess,0,from}

	tsIns = append(tsIns,tsInput)

	tsOutput := &TranslationOutput{int64(amount),to}

	tsOuts = append(tsOuts,tsOutput)

	ts := &Translation{[]byte{},tsIns,tsOuts}

	ts.HashTransaction()

	return ts
}

//判断交易是否是coinbase类型
func (ts *Translation)isCoinBaseTransaction() bool {
	return len(ts.TrsIns[0].TxHash) == 0 && ts.TrsOuts[0].Value == -1
}


