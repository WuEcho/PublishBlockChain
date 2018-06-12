package main

import (
	"code/0611/walet2/blc"
	"fmt"
)

func main()  {

	waltes := blc.NewWalets()

	waltes.CreatWalets()

	waltes.CreatWalets()

	fmt.Println(*waltes)
}