package blc

type UTXO struct {

	TXHash []byte

	Index int

	Output *TranslationOutput
} 