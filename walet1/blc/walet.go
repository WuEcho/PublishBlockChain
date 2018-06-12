package blc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"code/publicChain/part45-base58/BLC"
	"bytes"
)


type Walet struct {

	//私钥
	PrivateKey ecdsa.PrivateKey
	//公钥
	Publishkey []byte
}

const version = byte(0x00)
const addressChecksumLen = 4


func NewWalet() *Walet {
	privateKey,publishKey := newKeyPair()

	return &Walet{privateKey,publishKey}
}

//生成公钥以及私钥
func newKeyPair() (ecdsa.PrivateKey,[]byte) {
	curve := elliptic.P256()
	private,err := ecdsa.GenerateKey(curve,rand.Reader)
	if err != nil{
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(),private.PublicKey.Y.Bytes()...)
	return *private,pubKey
}

//通过计算地址的逆运算来验证地址是否合法
func IsValidOfAddress(address []byte) bool {

	version_publishkey_checknum := BLC.Base58Decode(address)

	checkNum := version_publishkey_checknum[len(version_publishkey_checknum)-addressChecksumLen:]

	version_publishkey := version_publishkey_checknum[:len(version_publishkey_checknum)-addressChecksumLen]

	hashed := CheckSum(version_publishkey)

	if bytes.Compare(hashed,checkNum) == 0 {
		return true
	}
	return false
}


//获取地址的方法
func (w *Walet)GetAddress() []byte {

	ripemd160 := w.Ripemd160Hash(w.Publishkey)

    version_ripemd := append([]byte{version},ripemd160...)

    checkByts := CheckSum(version_ripemd)

    bytes := append(version_ripemd,checkByts...)

    return Base58Encode(bytes)
}

//取前4字节
func CheckSum(payload []byte) []byte {

 	hash1 := sha256.Sum256(payload)

 	hash1 = sha256.Sum256(hash1[:])

 	return hash1[:addressChecksumLen]
}

//riemd160+sha256加密
func (w *Walet)Ripemd160Hash(publishKey []byte) []byte {
	//sha256
	hash256 := sha256.New()
	hash256.Write(publishKey)
	hash := hash256.Sum(nil)

	//ripemd160
	ripemd := ripemd160.New()
	ripemd.Write(hash)
	return ripemd.Sum(nil)
}