package main

import (
	"asng/proto"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	for i := 1; i < len(os.Args); i += 1 {
		hex, err := hex.DecodeString(os.Args[i])
		if err != nil {
			fmt.Printf("%s, decode err:%s\n", os.Args[i], err)
			continue
		}
		cursor := 0
		v, err := proto.ReadTDXFloat(hex, &cursor)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s: %f", os.Args[i], v)
	}
}
