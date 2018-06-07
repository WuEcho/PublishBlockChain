package main

import (
	"code/0605demo/BlockChainDemo1/blc"
	"fmt"
)

func main()  {
   blockChain := blc.NewGenesisBlockChain()

   blockChain.AddBlock("第二区块")

	for i := 0;i < len(blockChain.Block) ; i++ {
		block := blockChain.Block[i]
		fmt.Printf("时间戳： %d \n",block.TimeStamp)
		fmt.Printf("哈希值： %x \n",block.Hash)
		fmt.Printf("前一区块哈希： %x \n",block.PrefHash)
		fmt.Printf("区块信息： %s \n",block.Data)
		fmt.Printf("当前区块高度: %d \n",block.Height)
		fmt.Println("--------------------------")
	}

}