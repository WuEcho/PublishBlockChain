package blc

type TranslationOutput struct {
    //转账金额
	Value int64
    //目标账户
	ScriptPubKey string
}

//解锁
func (tsout *TranslationOutput)UnlockScriptPubKeyWithAddress(address string) bool {
	return tsout.ScriptPubKey == address
}