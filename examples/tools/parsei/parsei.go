package main

import (
	"asng/proto"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	for i := 1; i < len(os.Args); i += 1 {
		c := 0
		bytes, err := hex.DecodeString(os.Args[i])
		if err != nil {
			fmt.Printf("%s, decode string err:%s\n", os.Args[i], err)
			continue
		}
		v, err := proto.ReadTDXInt(bytes, &c)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s, %d\n", os.Args[i], v)
	}
}
