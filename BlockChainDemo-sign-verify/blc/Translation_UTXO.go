package blc
//UTXO
type Translation struct {
	//事务哈希
	TransHash  []byte
	//输入事务
	TrsIns     []*TranslationInput
	//输出事务
	TrsOuts    []*TranslationOutput
}