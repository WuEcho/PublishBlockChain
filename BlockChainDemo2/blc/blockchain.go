package blc

import (
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"fmt"
	"time"
)

type BlockChain struct {
	Tip   []byte   //当前区块哈希
	DB    *bolt.DB  //数据库对象
}
//数据库名
const databaseName = "blockChain"
//表名
const tableName = "block"

func NewGenesisBlockChain() *BlockChain {
	//原来的代码
	//return &BlockChain{[]*Block{NewGenesisBlock("创世区块")}}
	//将区块信息保存到数据库中

	block := NewGenesisBlock("创世块--")

    db,err := bolt.Open(databaseName,0600,nil)
	if err != nil {
		log.Panic(err)
	}

	var hashByte []byte

	db.Update(func(tx *bolt.Tx) error {
		//创建数据库表
		b,err := tx.CreateBucketIfNotExists([]byte(tableName))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(block.Hash,block.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"),block.Hash)
		if err != nil {
			log.Panic(err)
		}

		hashByte = block.Hash
       return err
	})
  return &BlockChain{hashByte,db}
}

func (blockChain *BlockChain)AddBlock(data string) {

	//prefBlock := blockChain.Block[len(blockChain.Block)-1]
	//
	//block := NewBlock(data,prefBlock.Height+1,prefBlock.Hash)
	//
	//blockChain.Block = append(blockChain.Block,block)

	//将新创建的区块存入数据库
	//1. 先从数据库中将当前区块取出
	err := blockChain.DB.Update(func(tx *bolt.Tx) error {

		//获取表
		b := tx.Bucket([]byte(tableName))
		if b != nil {
			blockData := b.Get([]byte(blockChain.Tip))
			//将区块反序列化
			block := DeserializeBlock(blockData)

			//创建新的区块
			newBlock := NewBlock(data,block.Height+1,block.Hash)
            //将区块存入数据库
            err := b.Put(newBlock.Hash,newBlock.Serialize())
            if err != nil{
            	log.Panic(err)
			}

			//更新数据库存储"l"对应的值
			err = b.Put([]byte("l"),newBlock.Hash)
			if err != nil{
				log.Panic(err)
			}

			//更新blockChain对象的Tig
			blockChain.Tip = newBlock.Hash
		}
		return nil
	})

	if err != nil{
		log.Panic(err)
	}
}

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


func (blockChain *BlockChain)PrintChain()  {
		////遍历数据库
		//err := blockChain.DB.View(func(tx *bolt.Tx) error {
		//
		//	b := tx.Bucket([]byte(tableName))
		//
		//	var bigInt  big.Int
		//
		//	if b != nil {
		//		currentHash := blockChain.Tip
		//		for  {
		//			blcokBytes := b.Get(currentHash)
		//
		//			block := DeserializeBlock(blcokBytes)
		//
		//			bigInt.SetBytes(currentHash)
		//
		//			if big.NewInt(0).Cmp(&bigInt) == 0{
		//				break
		//			}
		//			currentHash = block.PrefHash
		//
		//			fmt.Println(block)
		//
		//
		//			fmt.Printf("时间戳： %d \n",block.TimeStamp)
		//			fmt.Printf("哈希值： %x \n",block.Hash)
		//			fmt.Printf("前一区块哈希： %x \n",block.PrefHash)
		//			fmt.Printf("区块信息： %s \n",block.Data)
		//			fmt.Printf("当前区块高度: %d \n",block.Height)
		//			fmt.Println("--------------------------")
		//		}
		//	}
		//	return nil
		// })

	//if err != nil {
	//	log.Panic(err)
	//}

	var block *Block

	var bigInt big.Int
    //初始化迭代器
	blockchainIterator := blockChain.Iterator()

	for  {
		//遍历上一区块
		block = blockchainIterator.NextBlock()

		fmt.Println("------------------")
		fmt.Printf("Height：%d\n",block.Height)
		fmt.Printf("PrevBlockHash：%x\n",block.PrefHash)
		fmt.Printf("Data：%s\n",block.Data)
		//格式化时间
		fmt.Printf("Timestamp：%s\n",time.Unix(block.TimeStamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash：%x\n",block.Hash)
		fmt.Printf("Nonce：%d\n",block.Nonce)
        fmt.Println("------------------")

		bigInt.SetBytes(block.PrefHash)
		//当区块的前一区块值都为零的时候即遍历至创世区块
		if big.NewInt(0).Cmp(&bigInt) == 0 {
             break
		}
	}



}