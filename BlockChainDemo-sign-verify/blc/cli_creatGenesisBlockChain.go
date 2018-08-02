package blc


func (cli *Cli)creatGenesisBlock(address string)  {

	blockChain := NewGenesisBlockChain(address)
	defer blockChain.DB.Close()
}
