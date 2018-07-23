package blc

import (
	"bytes"
	"encoding/binary"
	"log"
	"encoding/json"
	"fmt"
)

func IntToHexo(number int64) []byte {
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer,binary.BigEndian,number)

	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

//json转换[]string
func JSONToArray(jsonString string) []string {

	//json转换[]string
	var sArr  []string
	if err := json.Unmarshal([]byte(jsonString),&sArr); err != nil{
		log.Panic(err)
		fmt.Println("error:",err)
	}
	return sArr
}