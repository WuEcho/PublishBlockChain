package blc

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	//前一区块的哈希
	PrefHash []byte
	//交易信息
	Data     []byte
    //当前区块哈希
	Hash     []byte
    //当前区块高度
	Height    int64
    //随机值
	Nonce     int64
    //时间戳
	TimeStamp int64
}

func NewBlock(data string,height int64,prefHash []byte) *Block {
   block := &Block{
   	Data:[]byte(data),
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

func NewGenesisBlock(data string) *Block {
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