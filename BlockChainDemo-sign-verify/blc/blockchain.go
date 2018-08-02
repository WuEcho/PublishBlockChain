package blc

import (
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"fmt"
	"time"
	"os"
	"strconv"
	"encoding/hex"
	"crypto/ecdsa"
	"bytes"
)

type BlockChain struct {
	Tip   []byte   //当前区块哈希
	DB    *bolt.DB  //数据库对象
}
//数据库名
const databaseName = "blockChain"
//表名
const tableName = "block"

func NewGenesisBlockChain(address string) *BlockChain {

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

			//创建coinbases交易
			tscb := NewCoinBaseTransaction(address)


			block = NewGenesisBlock([]*Translation{tscb})

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

func (blockChain *BlockChain)SignTransaction(tx *Translation,privKey ecdsa.PrivateKey)  {

	if tx.isCoinBaseTransaction() {
		return
	}

	prevTXs := make(map[string]Translation)
	for _,vin := range tx.TrsIns {
		prevTx,err := blockChain.FindTransaction(vin.TxHash)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTx.TransHash)] = prevTx
	}

	tx.Sign(privKey,prevTXs)
}

func (blockChain *BlockChain)FindTransaction(id []byte) (Translation,error) {
	bci := blockChain.Iterator()

	for {
		block := bci.NextBlock()

		for _,tx := range block.Translations {
			if bytes.Compare(tx.TransHash,id) == 0 {
				return *tx,nil
			}
		}
		var hashInt big.Int
		hashInt.SetBytes(block.PrefHash)

		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
	return Translation{},nil
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


//挖取新的区块
func (blockChain *BlockChain)MineNewBlock(from []string,to []string,amount []string) {

	var tss  []*Translation
	var block  *Block

	blockChain.DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(tableName))
		if b != nil {
		   hash := b.Get([]byte("l"))

		   blockBytes := b.Get(hash)

		   block = DeserializeBlock(blockBytes)
		}

		return nil
	})

	for index,address := range from {
		value,_ := strconv.Atoi(amount[index])
		tx := NewSimpleTranslation(address,to[index],value,blockChain,tss)
		tss = append(tss,tx)
	}

	//在新建区块之前对数字签名进行验证
	for _,tx := range tss {
		if blockChain.VerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}


	//新建区块
	block = NewBlock(tss,block.Height+1,block.Hash)
	blockChain.DB.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(tableName))

		if b != nil {
			b.Put(block.Hash,block.Serialize())
			b.Put([]byte("l"),block.Hash)
			blockChain.Tip = block.Hash
		}
		return nil
	})

}

//验证数字签名
func (blockChain *BlockChain)VerifyTransaction(tx *Translation) bool {
	prevTxs := make(map[string]Translation)

	for _,vin := range tx.TrsIns{
		prevTx,err := blockChain.FindTransaction(vin.TxHash)
		if err != nil {
			log.Panic(err)
		}
		prevTxs[hex.EncodeToString(prevTx.TransHash)] = prevTx
	}
	return tx.Verify(prevTxs)
}


//如果一个地址对应的TrsOuts未花费，那么这个Translation添加到数组中
func (blockChain *BlockChain)UnUTXOs(address string,txs []*Translation) []*UTXO {

	var unUTXOs []*UTXO

	spentTXOutputs := make(map[string][]int)

	//遍历交易数组中的未打包的交易
	for _,tx := range txs {
		if tx.isCoinBaseTransaction() == false {
			for _,in := range tx.TrsIns {
				//是否能解锁
				publishKey := Base58Decode([]byte(address))

				ripemd160Hash := publishKey[1:len(publishKey)-4]

				if in.UnlockRipemd160Hash(ripemd160Hash) {
					key := hex.EncodeToString(in.TxHash)
					spentTXOutputs[key] = append(spentTXOutputs[key],in.VoutInde)
				}
			}
		}
	}

	for _,tx := range txs {
      label:
		for index,txout := range tx.TrsOuts {
			if txout.UnlockScriptPubKeyWithAddress(address) {
				if len(spentTXOutputs) == 0 {
					utxo := &UTXO{TXHash:tx.TransHash,Index:index,Output:txout}
					unUTXOs = append(unUTXOs,utxo)
				}else {
					for hash,indexArray := range spentTXOutputs {
						txHashStr := hex.EncodeToString(tx.TransHash)
						if hash == txHashStr {
							var isUnSpentUTXO bool = false
							for _,outIndex := range indexArray {
								if outIndex == index {
									isUnSpentUTXO = true
									continue label
								}
								if isUnSpentUTXO ==false {
									utxo := &UTXO{tx.TransHash,index,txout}
									unUTXOs = append(unUTXOs,utxo)
								}
							}
						}else {
							utxo := &UTXO{tx.TransHash,index,txout}
							unUTXOs = append(unUTXOs,utxo)
						}
					}
				}

			}
		}
	}

	blockIterator := blockChain.Iterator()

	for {
		block := blockIterator.NextBlock()
		for _,ts := range block.Translations {

			if ts.isCoinBaseTransaction() == false {
				for _,in := range ts.TrsIns {
					//是否能够解锁

					publishHash := Base58Decode([]byte(address))

					ripemd160Hash := publishHash[1:len(publishHash)-4]


					if in.UnlockRipemd160Hash(ripemd160Hash) {

						key := hex.EncodeToString(in.TxHash)

                        spentTXOutputs[key] = append(spentTXOutputs[key],in.VoutInde)
					}
				}

			}
			
		work:
			for index,out := range ts.TrsOuts {

				if out.UnlockScriptPubKeyWithAddress(address) {

					if spentTXOutputs != nil {
						if len(spentTXOutputs) != 0 {

							var isSpentUTXO bool = false
							
							for tsHash,indexArray := range spentTXOutputs {

								for _,i := range indexArray {

									if tsHash== hex.EncodeToString(ts.TransHash) && index == i {
										isSpentUTXO = true
										continue work
									}
								}
							}
							if isSpentUTXO==false {
							  utxo := &UTXO{ts.TransHash,index,out}
							  unUTXOs = append(unUTXOs,utxo)
							}

						}else {
							utxo := &UTXO{ts.TransHash,index,out}
							unUTXOs = append(unUTXOs,utxo)
						}
					}

				}

			}

		}
		var hashInt big.Int
		hashInt.SetBytes(block.PrefHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}

	}
  return unUTXOs
}

func (blockchain *BlockChain)FindSpendableUTXOs(from string,amount int,txs []*Translation) (int64,map[string][]int) {

	utxos := blockchain.UnUTXOs(from,txs)

	spendabelUTXO := make(map[string][]int)

	var value  int64
	for _,utxo := range utxos{
       value = value + utxo.Output.Value
       hash := hex.EncodeToString(utxo.TXHash)
       spendabelUTXO[hash] = append(spendabelUTXO[hash],utxo.Index)

		if value >= int64(amount) {
			break
		}
	}

	fmt.Printf("余额为%d----需要支付的金额为：%d \n",value,int64(amount))

	if value < int64(amount) {
		fmt.Printf("%s's found is 不足\n", from)
		os.Exit(1)
	}
	return value,spendabelUTXO
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
		fmt.Printf("Data：%v\n",block.Translations)
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

func (blockchain *BlockChain)GetBalance(address string) int64 {
	utxos := blockchain.UnUTXOs(address,[]*Translation{})

	var amount int64

	for _, utxo := range utxos{
		amount = amount + utxo.Output.Value
	}
	return amount
}