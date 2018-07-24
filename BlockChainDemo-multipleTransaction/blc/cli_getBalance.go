package blc

import "fmt"

//查询余额
func (cli *Cli)getBalance(address string)  {
	fmt.Println("地址："+address)

	blockChain := BlockChainObject()
	defer blockChain.DB.Close()

	amount := blockChain.GetBalance(address)

	fmt.Printf("%s 一共有%d个token",address,amount)
}