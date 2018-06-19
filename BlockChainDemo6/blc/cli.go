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

	blockFlag := flag.NewFlagSet("addBlock",flag.ExitOnError)

	printFlag := flag.NewFlagSet("printBlock",flag.ExitOnError)

	creatGenesisBlock := flag.NewFlagSet("creatGenesisBlock",flag.ExitOnError)

	//添加指令方法
	sendFlag := flag.NewFlagSet("send",flag.ExitOnError)

	blockChainFlag := blockFlag.String("data","添加区块","添加新的block")

	genesisBlockFlag := creatGenesisBlock.String("data","创世区块","创建创世区块")

	flagFrom := sendFlag.String("from","","转账源地址")

	flagTo := sendFlag.String("to","","转账目标地址")

	flagAmount := sendFlag.String("amount","","转账金额")

	switch os.Args[1] {

	case "creatGenesisBlock":
		err := creatGenesisBlock.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "addBlock":
		err := blockFlag.Parse(os.Args[2:])
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

	if blockFlag.Parsed() {
		if *blockChainFlag == "" {
			printUsage()
			os.Exit(0)
		}

		//添加区块
        fmt.Println("添加区块")
		cli.addBlock([]*Translation{})
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

}

func (cli *Cli)creatGenesisBlock(address string)  {
	 NewGenesisBlockChain(address)
}


func (cli *Cli)addBlock(data []*Translation) {

     blockChain := BlockChainObject()
     defer blockChain.DB.Close()

     blockChain.AddBlock(data)
}

func (cli *Cli)printChain()  {
	blockChain := BlockChainObject()
	defer blockChain.DB.Close()

	blockChain.PrintChain()
}

func MineNewBlock(fromString []string,toString []string,amountString []string)  {

   blc := BlockChainObject()
   defer blc.DB.Close()

   blc.MineNewBlock(fromString,toString,amountString)
}


func printUsage()  {
	fmt.Println("Usage----")
	//添加创建发送交易的方法
	fmt.Println("\t creatGenesisBlock  --- 创建创世区块")
	fmt.Println("\t send -from -to -amount --- 添加新的区块")
	fmt.Println("\t addBlock -data --- 添加新的区块")
	fmt.Println("\t printBlock  --- 打印区块数据")
}