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
	//新增
	//添加创建创世区块的指令
	creatGenesisBlock := flag.NewFlagSet("creatGenesisBlock",flag.ExitOnError)

	blockChainFlag := blockFlag.String("data","添加区块","添加新的block")

	genesisBlockFlag := creatGenesisBlock.String("data","创世区块","创建创世区块")

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
		cli.addBlock(*blockChainFlag)
	}

	if printFlag.Parsed() {

		fmt.Println(cli.Blc)

		//打印数据
		cli.printChain()
	}
}

func (cli *Cli)creatGenesisBlock(data string)  {
	 NewGenesisBlockChain(data)
}


func (cli *Cli)addBlock(data string) {

     blockChain := BlockChainObject()
     defer blockChain.DB.Close()

     blockChain.AddBlock(data)
}

func (cli *Cli)printChain()  {
	blockChain := BlockChainObject()
	defer blockChain.DB.Close()

	blockChain.PrintChain()
}


func printUsage()  {
	fmt.Println("Usage----")
	fmt.Println("\t creatGenesisBlock  --- 创建创世区块")
	fmt.Println("\t addBlock -data --- 添加新的区块")
	fmt.Println("\t printBlock  --- 打印区块数据")
}