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

	blockChainFlag := blockFlag.String("data","创世区块","添加新的block")

	switch os.Args[1] {
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

	if blockFlag.Parsed() {
		if *blockChainFlag == "" {
			printUsage()
			os.Exit(0)
		}

		//添加区块
        fmt.Println("添加区块")
		cli.Blc.AddBlock(*blockChainFlag)
	}

	if printFlag.Parsed() {
		//打印数据
		cli.Blc.PrintChain()
	}
}

func printUsage()  {
	fmt.Println("Usage----")
	fmt.Println("\t addBlock -data --- 添加新的区块")
	fmt.Println("\t printBlock  --- 打印区块数据")

}