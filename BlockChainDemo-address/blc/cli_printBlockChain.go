package blc


func (cli *Cli)printChain()  {
	blockChain := BlockChainObject()
	defer blockChain.DB.Close()

	blockChain.PrintChain()
}


