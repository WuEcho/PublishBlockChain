package main

import (
	"code/0605demo/BlockChainDemo2/blc"
)

func main()  {
   blockChain := blc.NewGenesisBlockChain()

   blockChain.AddBlock("第二区块")

   blockChain.PrintChain()
}