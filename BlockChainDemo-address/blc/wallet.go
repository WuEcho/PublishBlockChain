package blc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"log"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"fmt"
	"bytes"
)

const version = byte(0x00)
const addressChecksumLen = 4

type Wallet struct {
	//1. 私钥
	PrivateKey ecdsa.PrivateKey

	//2. 公钥
	PublicKey  []byte
}

func IsValidForAdress(adress []byte) bool {

	version_public_checksumBytes := Base58Decode(adress)

	fmt.Println(version_public_checksumBytes)

	//25
	//4
	//21
	checkSumBytes := version_public_checksumBytes[len(version_public_checksumBytes) - addressChecksumLen:]

	version_ripemd160 := version_public_checksumBytes[:len(version_public_checksumBytes) - addressChecksumLen]

	fmt.Println(checkSumBytes)
	fmt.Println(version_ripemd160)

	checkBytes := CheckSum(version_ripemd160)
	fmt.Println(checkBytes)

	if bytes.Compare(checkSumBytes,checkBytes) == 0 {
		return true
	}

	return false
}


func (w *Wallet) GetAddress() []byte  {

	//1. hash160

	ripemd160Hash := w.Ripemd160Hash(w.PublicKey)

	version_ripemd160Hash := append([]byte{version},ripemd160Hash...)

	checkSumBytes := CheckSum(version_ripemd160Hash)

	bytes := append(version_ripemd160Hash,checkSumBytes...)

	return Base58Encode(bytes)
}

func CheckSum(payload []byte) []byte {

	fmt.Println(payload)

	hash1 := sha256.Sum256(payload)

	fmt.Println(hash1)

	hash2 := sha256.Sum256(hash1[:])

	fmt.Println(hash2)
	return hash2[:addressChecksumLen]
}


func (w *Wallet) Ripemd160Hash(publicKey []byte) []byte {

	//1. 256
	hash256 := sha256.New()
	hash256.Write(publicKey)
	hash := hash256.Sum(nil)

	//2. 160
	ripemd160 := ripemd160.New()
	ripemd160.Write(hash)

	return ripemd160.Sum(nil)
}

// 创建钱包
func NewWallet() *Wallet {

	privateKey,publicKey := newKeyPair()

	return &Wallet{privateKey,publicKey}
}


// 通过私钥产生公钥
func newKeyPair() (ecdsa.PrivateKey,[]byte) {

	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}