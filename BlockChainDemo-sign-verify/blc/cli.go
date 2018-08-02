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

	//查询余额
	getBalanceFlag := flag.NewFlagSet("getBalance",flag.ExitOnError)
	getBalanceWithAddress := getBalanceFlag.String("address","","查询目标地址余额")


	//交易
	sendFlag := flag.NewFlagSet("send",flag.ExitOnError)

	genesisBlockFlag := creatGenesisBlock.String("data","创世区块","创建创世区块")

	flagFrom := sendFlag.String("from","","转账源地址")

	flagTo := sendFlag.String("to","","转账目标地址")

	flagAmount := sendFlag.String("amount","","转账金额")


	//创建钱包
	creatWallet := flag.NewFlagSet("creatWallet",flag.ExitOnError)
	//遍历钱包
	addressList := flag.NewFlagSet("addressLists",flag.ExitOnError)
	

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
	case "creatWallet":
		err := creatWallet.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "addressList":
		err := addressList.Parse(os.Args[2:])
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
		}else if IsValidForAdress([]byte(*genesisBlockFlag))==false{
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

		from := JSONToArray(*flagFrom)
		to := JSONToArray(*flagTo)

		for index,address := range from {
			if IsValidForAdress([]byte(address))==false || IsValidForAdress([]byte(to[index])) ==false{
				fmt.Printf("地址无效......")
				printUsage()
				os.Exit(1)
			}
		}

		//将参数传递进入创建新的区块的方法
		MineNewBlock(from,to,JSONToArray(*flagAmount))
	}

	if getBalanceFlag.Parsed() {

		if *getBalanceWithAddress == ""{
			fmt.Println("地址不能为空")
			printUsage()
			os.Exit(1)
		}else if (IsValidForAdress([]byte(*getBalanceWithAddress))==false){
			fmt.Println("地址无效....")
			printUsage()
			os.Exit(1)
		}

		cli.getBalance(*getBalanceWithAddress)
	}

	if creatWallet.Parsed() {
		cli.creatWallet()
	}

	if addressList.Parsed() {
		cli.addressList()
	}
}


func printUsage()  {
	fmt.Println("Usage----")
	fmt.Println("\t creatWallet   --- 创建钱包地址")
	fmt.Println("\t addressList   --- 输出所有钱包地址")
	fmt.Println("\t creatGenesisBlock  -data --- 创建创世区块")
	fmt.Println("\t send -from -to -amount --- 添加新的区块")
	fmt.Println("\t getBalance -address  --- 查询余额")
	fmt.Println("\t printBlock  --- 打印区块数据")
}