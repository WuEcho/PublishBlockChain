package blc

import "time"

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

