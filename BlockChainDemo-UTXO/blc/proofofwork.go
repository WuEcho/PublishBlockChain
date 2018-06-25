package blc

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

type ProofOfWork struct {
	Block *Block
	targetNumber *big.Int
}

const targeNumber = 16

func NewProofofWork(block *Block) *ProofOfWork {

	 target := big.NewInt(1)

	 target = target.Lsh(target,256-targeNumber)

	return &ProofOfWork{block,target}
}

func (pow *ProofOfWork)Run() ([]byte,int64) {
	var bigInt big.Int

	var hash [32]byte

	nonce := 0

	for {
        bytesHash := pow.prepareData(int64(nonce))

        hash = sha256.Sum256(bytesHash)

        fmt.Printf("\r ----%x \n",hash)

        bigInt.SetBytes(hash[:])

		if pow.targetNumber.Cmp(&bigInt) == 1 {
			break
		}
		nonce = nonce +1
	}
	return hash[:],int64(nonce)
}

func (pow *ProofOfWork)prepareData(nonce int64) []byte {
	hashByte := bytes.Join([][]byte{
		pow.Block.PrefHash,
		pow.Block.HashTranslation(),
		IntToHexo(nonce),
		IntToHexo(pow.Block.Height),
		IntToHexo(pow.Block.TimeStamp),
	},[]byte{})

	return hashByte
}


func (p *ProofOfWork)IsValid() bool {

	var hashInt *big.Int

	hashInt.SetBytes(p.Block.Hash)

	if p.targetNumber.Cmp(hashInt) ==1 {
		return true
	}
    return false
}