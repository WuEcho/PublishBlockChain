package blc

import "github.com/boltdb/bolt"


//迭代结构体
type BlockchainIterator struct {
	//当前区块哈希
	CurrentHash []byte
	//数据库
	Db   *bolt.DB
}

//实例化
func (blcokChain *BlockChain)Iterator() *BlockchainIterator {
	return &BlockchainIterator{blcokChain.Tip,blcokChain.DB}
}


func (blockchainIterator *BlockchainIterator)NextBlock() *Block {

	var block  *Block

	blockchainIterator.Db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(tableName))
		if b != nil {

			bts := b.Get(blockchainIterator.CurrentHash)

			block = DeserializeBlock(bts)

			blockchainIterator.CurrentHash = block.PrefHash
		}

		return nil
	})


	return block
}

