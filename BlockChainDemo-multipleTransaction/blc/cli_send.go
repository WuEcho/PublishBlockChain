package blc

import "fmt"

func MineNewBlock(fromString []string,toString []string,amountString []string)  {

	blc := BlockChainObject()
	defer blc.DB.Close()

	fmt.Println()

	blc.MineNewBlock(fromString,toString,amountString)
}