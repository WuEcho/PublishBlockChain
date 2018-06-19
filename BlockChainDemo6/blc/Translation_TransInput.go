package blc

type TranslationInput struct {
	//事务的哈希
	TxHash      []byte
	//在trsout中的序号
	VoutInde    int
	//地址
	Address     string
}
