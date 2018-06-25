package blc

import (
	"flag"
	"os"
	"fmt"
	"log"
)

type Cli struct {
	Blc *BlockChain
}

//判断终端读取到的内容合法性
func isVailid()  {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}
}


func (cli *Cli)Run()  {

	isVailid()

	printFlag := flag.NewFlagSet("printBlock",flag.ExitOnError)

	creatGenesisBlock := flag.NewFlagSet("creatGenesisBlock",flag.ExitOnError)

	getBalanceFlag := flag.NewFlagSet("getBalance",flag.ExitOnError)

	//添加指令方法
	sendFlag := flag.NewFlagSet("send",flag.ExitOnError)

	genesisBlockFlag := creatGenesisBlock.String("data","创世区块","创建创世区块")

	flagFrom := sendFlag.String("from","","转账源地址")

	flagTo := sendFlag.String("to","","转账目标地址")

	flagAmount := sendFlag.String("amount","","转账金额")

	getBalanceWithAddress := getBalanceFlag.String("address","","查询目标地址余额")


	switch os.Args[1] {

	case "creatGenesisBlock":
		err := creatGenesisBlock.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "printBlock":
		err := printFlag.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "send":
		err := sendFlag.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}

	case "getBalance":
		err := getBalanceFlag.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(0)
	}

	if creatGenesisBlock.Parsed() {
		if *genesisBlockFlag == "" {
			printUsage()
			os.Exit(1)
		}
		cli.creatGenesisBlock(*genesisBlockFlag)
	}

	if printFlag.Parsed() {

		fmt.Println(cli.Blc)

		//打印数据
		cli.printChain()
	}

	if sendFlag.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == ""{
			printUsage()
			os.Exit(1)
		}

		//将参数传递进入创建新的区块的方法
		MineNewBlock(JSONToArray(*flagFrom),JSONToArray(*flagTo),JSONToArray(*flagAmount))
	}

	if getBalanceFlag.Parsed() {

		if *getBalanceWithAddress == ""{
			fmt.Println("地址不能为空")
			printUsage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceWithAddress)
	}

}

func (cli *Cli)creatGenesisBlock(address string)  {
	 NewGenesisBlockChain(address)
}


func (cli *Cli)printChain()  {
	blockChain := BlockChainObject()
	defer blockChain.DB.Close()

	blockChain.PrintChain()
}

//查询余额
func (cli *Cli)getBalance(address string)  {
	fmt.Println("地址："+address)

   blockChain := BlockChainObject()
   defer blockChain.DB.Close()

   tsOutPut := blockChain.UnUTXOs(address)

   fmt.Println("========")

   for _,out := range tsOutPut{
   	fmt.Println(out)
   }
}

func MineNewBlock(fromString []string,toString []string,amountString []string)  {

   blc := BlockChainObject()
   defer blc.DB.Close()

   blc.MineNewBlock(fromString,toString,amountString)
}


func printUsage()  {
	fmt.Println("Usage----")
	fmt.Println("\t creatGenesisBlock  --- 创建创世区块")
	fmt.Println("\t send -from -to -amount --- 添加新的区块")
	fmt.Println("\t getBalance -address  --- 查询余额")
	fmt.Println("\t printBlock  --- 打印区块数据")
}