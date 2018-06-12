package main

import (
	"code/0611/walet1/blc"
	"fmt"
)

func main()  {

	walet := blc.NewWalet()

	waltAddress := walet.GetAddress()

	fmt.Println(string(waltAddress))

	addressIsValid := blc.IsValidOfAddress(waltAddress)

	string := ""

	if addressIsValid {
		string = "是"
	}else {
		string = "否"
	}

	fmt.Printf("地址是否合法:%s",string)
}