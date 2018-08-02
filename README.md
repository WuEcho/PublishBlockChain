# PublishBlockChain
Go语言实现简易公链

博客Demo示例解析地址：http://www.wuecho.com/2018/06/05/Go语言实现搭建公链-1/
                    http://www.wuecho.com/2018/06/11/Go语言实现公链-2-地址与钱包/
                    http://www.wuecho.com/2018/06/13/Go语言实现公链-3-事务-交易/


针对BlockChainDemo需要复制到Gopath里面，修改头文件的引入，才能正常运行。

对于BlockChainDemo1，以及BlockChainDemo2正常复制到gopath目录下运行即可。
对于BlockChainDemo3，以及其他需要打开命令行使用终端命令：

&cd BlockChainDemo3     //进入到BlockChainDemo3文件夹下，或者进入到BlockChainDemo4
&go build main.go      //编译main.go
//使用以下指令控制代码
&./main creatGenesisBlock 
&./main addBlock -data 
&./main printBlock 






