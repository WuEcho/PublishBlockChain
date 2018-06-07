package blc


type BlockChain struct {
	Block []*Block
}


func NewGenesisBlockChain() *BlockChain {
	return &BlockChain{[]*Block{NewGenesisBlock("创世区块")}}
}

func (blockChain *BlockChain)AddBlock(data string) {

	prefBlock := blockChain.Block[len(blockChain.Block)-1]

	block := NewBlock(data,prefBlock.Height+1,prefBlock.Hash)

	blockChain.Block = append(blockChain.Block,block)
}