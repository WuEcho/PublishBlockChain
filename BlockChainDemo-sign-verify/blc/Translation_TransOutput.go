package blc

import "bytes"

type TranslationOutput struct {
    //转账金额
	Value int64
    //目标账户
	Ripemd160Hash []byte
}

//解锁
func (tsout *TranslationOutput)UnlockScriptPubKeyWithAddress(address string) bool {

	publicKeyHash := Base58Decode([]byte(address))

	hash160 := publicKeyHash[1:len(publicKeyHash)-4]

	return bytes.Compare(tsout.Ripemd160Hash,hash160) == 0
}

//锁定
func (tsout *TranslationOutput)Lock(address string)  {
	publicKeyHash := Base58Decode([]byte(address))

	tsout.Ripemd160Hash = publicKeyHash[1:len(publicKeyHash)-4]
}


func NewTsOutput(value int64,address string) *TranslationOutput {

	txout := &TranslationOutput{value,nil}

	txout.Lock(address)

	return txout
}