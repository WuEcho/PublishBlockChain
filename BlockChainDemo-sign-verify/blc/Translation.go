package blc

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/elliptic"
	"math/big"
)



//创建CoinBase类型交易
func NewCoinBaseTransaction(address string) *Translation {

  tsInput := &TranslationInput{[]byte{},-1,nil,[]byte{}}

  tsOutput := NewTsOutput(10,address)

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


func NewSimpleTranslation(from string,to string,amount int,blockchain *BlockChain,txs []*Translation) *Translation {

	wallets,_ := NewWallets()

	wallet := wallets.WalletsMap[from]

	var tsIns []*TranslationInput
	var tsOuts  []*TranslationOutput

	money,spendableUTXODic := blockchain.FindSpendableUTXOs(from,amount,txs)

	if money < int64(amount) {
		log.Panic(errors.New("转账金额不足"))
		return nil
	}

	//只有转账金额大于
	for tsHash,indexArray := range spendableUTXODic {
         tsHashBytes,_ := hex.DecodeString(tsHash)
		for _,index := range indexArray{
			tsInput := &TranslationInput{tsHashBytes,index,nil,wallet.PublicKey}
			tsIns = append(tsIns,tsInput)
		}
	}

  //转账
  tsout := NewTsOutput(int64(amount),to)
  tsOuts = append(tsOuts,tsout)

  //找零
  tsout = NewTsOutput(int64(money)-int64(amount),from)
  tsOuts = append(tsOuts,tsout)

  ts := &Translation{[]byte{},tsIns,tsOuts}
  ts.HashTransaction()

  blockchain.SignTransaction(ts,wallet.PrivateKey)

  return ts
}

//判断交易是否是coinbase类型
func (ts *Translation)isCoinBaseTransaction() bool {
	return len(ts.TrsIns[0].TxHash) == 0 && ts.TrsOuts[0].Value == -1
}

func (ts *Translation)Hash() []byte  {
	txCopy := ts
	txCopy.TransHash = []byte{}

	hash := sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

func (ts *Translation)Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(ts)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

//数字签名
func (tx *Translation)Sign(privKey ecdsa.PrivateKey,prevTXs map[string]Translation)  {
	if tx.isCoinBaseTransaction() {
		return
	}

	for _,vins := range tx.TrsIns {
		if prevTXs[hex.EncodeToString(vins.TxHash)].TransHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmCopy()

	for inId , vin := range txCopy.TrsIns {
		prevTx := prevTXs[hex.EncodeToString(vin.TxHash)]
		txCopy.TrsIns[inId].Signature = nil
		txCopy.TrsIns[inId].PublicKey = prevTx.TrsOuts[vin.VoutInde].Ripemd160Hash
		txCopy.TransHash = txCopy.Hash()
		txCopy.TrsIns[inId].PublicKey = nil

		//签名
		r,s,err := ecdsa.Sign(rand.Reader,&privKey,txCopy.TransHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(),s.Bytes()...)

		tx.TrsIns[inId].Signature = signature
	}

}

//拷贝一份新的Transaction用于签名
func (tx *Translation)TrimmCopy() Translation {
	var inputs []*TranslationInput
	var outputs []*TranslationOutput

	for _,vin := range tx.TrsIns{
		inputs = append(inputs,&TranslationInput{vin.TxHash,vin.VoutInde,nil,nil})
	}

	for _,vout := range tx.TrsOuts{
		outputs = append(outputs,&TranslationOutput{vout.Value,vout.Ripemd160Hash})
	}

	txCopy := Translation{tx.TransHash,inputs,outputs}

	return txCopy
}

//验证数字签名
func (tx *Translation)Verify(prevTXs map[string]Translation) bool {
	if tx.isCoinBaseTransaction() {
		return true
	}

	for _,vin := range tx.TrsIns{
		if prevTXs[hex.EncodeToString(vin.TxHash)].TransHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmCopy()
	curve := elliptic.P256()

	for inID,vin := range tx.TrsIns{
		prevTx := prevTXs[hex.EncodeToString(vin.TxHash)]
		txCopy.TrsIns[inID].Signature = nil
		txCopy.TrsIns[inID].PublicKey = prevTx.TrsOuts[vin.VoutInde].Ripemd160Hash
		txCopy.TransHash = txCopy.Hash()
		txCopy.TrsIns[inID].PublicKey = nil

		// 私钥 ID
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen/2)])
		s.SetBytes(vin.Signature[(sigLen/2):])

		x := big.Int{}
		y := big.Int{}
 		keyLen := len(vin.PublicKey)
 		x.SetBytes(vin.PublicKey[:(keyLen/2)])
 		y.SetBytes(vin.PublicKey[(keyLen/2):])

 		rawPubKey := ecdsa.PublicKey{curve,&x,&y}

		if ecdsa.Verify(&rawPubKey,txCopy.TransHash,&r,&s) == false {
			return false
		}
	}
	return true
}
