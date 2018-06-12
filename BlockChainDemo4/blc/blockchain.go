package blc

import (
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"fmt"
	"time"
	"os"
)

type BlockChain struct {
	Tip   []byte   //当前区块哈希
	DB    *bolt.DB  //数据库对象
}
//数据库名
const databaseName = "blockChain"
//表名
const tableName = "block"

func NewGenesisBlockChain(data string) *BlockChain {

	var hashByte []byte

	var block *Block
	//判断数据库是否存在，如果已经存在那么说明数据库里面已经存有创世区块，不存在就保存
	if dbIsExist() {
		db,err := bolt.Open(databaseName,0600,nil)
		if err != nil {
			log.Panic(err)
		}

		db.Update(func(tx *bolt.Tx) error {

			//创建数据库表
			b, err := tx.CreateBucketIfNotExists([]byte(tableName))
			if err != nil {
				log.Panic(err)
			}
			block = NewGenesisBlock(data)

			err = b.Put(block.Hash, block.Serialize())

			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("l"), block.Hash)
			if err != nil {
				log.Panic(err)
			}

			return nil
		})
		return &BlockChain{hashByte,db}
	}

	blockChain := BlockChainObject()


  return blockChain
}

func BlockChainObject() *BlockChain {
	var block *Block

	//存在就从数据库中读取
	db,err := bolt.Open(databaseName,0600,nil)
	if err != nil {
		log.Panic(err)
	}

	db.View(func(tx *bolt.Tx) error {

		b :=  tx.Bucket([]byte(tableName))

		//存在就取出来
		bytesHash := b.Get([]byte("l"))

		blockBytes := b.Get(bytesHash)

		block = DeserializeBlock(blockBytes)

		return nil
	})
	return &BlockChain{block.Hash,db}
}



func dbIsExist() bool {
	if _,err := os.Stat(databaseName);os.IsNotExist(err) {
		return true
	}
	return false
}

func (blockChain *BlockChain)AddBlock(data string) {

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


func (blockChain *BlockChain)PrintChain()  {

	if dbIsExist() {
		fmt.Println("暂无数据")
		os.Exit(1)
	}

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