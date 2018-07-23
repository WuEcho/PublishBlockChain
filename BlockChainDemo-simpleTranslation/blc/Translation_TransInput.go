package blc

type TranslationInput struct {
	//事务的哈希
	TxHash      []byte
	//在trsout中的序号
	VoutInde    int
	//地址
	Address     string
}

//判断当前的消费的是谁的钱
func (in *TranslationInput)UnLockWithAdress(address string) bool {
	return in.Address == address
}
