package blc

import "fmt"

func (cli *Cli)creatWallet()  {

	wallets,_ := NewWallets()

	wallets.CreateNewWallet()

	fmt.Println(wallets.WalletsMap)
}
