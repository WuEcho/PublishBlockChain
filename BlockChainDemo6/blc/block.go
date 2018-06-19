package blc

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
)

type Block struct {
	//前一区块的哈希
	PrefHash []byte
	//交易信息
	Translation  []*Translation
    //当前区块哈希
	Hash     []byte
    //当前区块高度
	Height    int64
    //随机值
	Nonce     int64
    //时间戳
	TimeStamp int64
}



func NewBlock(data []*Translation,height int64,prefHash []byte) *Block {
   block := &Block{
   	Translation: data,
   	Height:height,
   	Nonce:0,
   	PrefHash:prefHash,
   	TimeStamp:time.Now().Unix(),
   }

   pow := NewProofofWork(block)

   hash,nonce := pow.Run()

   block.Hash = hash

   block.Nonce = nonce

   return block
}

func (block *Block)HashTranslation() []byte {

	var txHashes [][]byte
	var txhash  [32]byte
	for _,txs := range block.Translation{
		txHashes = append(txHashes,txs.TransHash)
	}
	txhash = sha256.Sum256(bytes.Join(txHashes,[]byte{}))

	return txhash[:]
}


func NewGenesisBlock(data []*Translation) *Block {
	return NewBlock(data,1,[]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}

// 将区块序列化成字节数组
func (block *Block) Serialize() []byte {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

//反序列化
func DeserializeBlock(blockBytes []byte) *Block {
	var block  Block
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Println(err)
	}
	return &block
}