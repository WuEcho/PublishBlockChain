package main

import (
	"code/0605demo/BlockChainDemo3/blc"
)

func main()  {

	blockChain := blc.NewGenesisBlockChain()

	cli := blc.Cli{blockChain}

	cli.Run()
}

