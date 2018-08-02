package blc

import "bytes"

type TranslationInput struct {
	//事务的哈希
	TxHash      []byte
	//在trsout中的序号
	VoutInde    int
	////地址
	//Address     string

	Signature []byte //数字签名

	PublicKey []byte //公钥
}


func (in *TranslationInput)UnlockRipemd160Hash(ripemd160Hash []byte) bool {

	publickey := Ripemd160Hash(in.PublicKey)

	return bytes.Compare(publickey,ripemd160Hash) == 0
}