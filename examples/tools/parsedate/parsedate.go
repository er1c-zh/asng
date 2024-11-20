package main

import (
	"bytes"
	"encoding/binary"
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

		var zipday, tminutes uint16
		binary.Read(bytes.NewBuffer(hex[0:2]), binary.LittleEndian, &zipday)
		binary.Read(bytes.NewBuffer(hex[2:4]), binary.LittleEndian, &tminutes)

		year := int((zipday >> 11) + 2004)
		month := int((zipday % 2048) / 100)
		day := int((zipday % 2048) % 100)
		hour := int(tminutes / 60)
		minute := int(tminutes % 60)

		fmt.Printf("%s: %d-%d-%d %d:%d\n", os.Args[i], year, month, day, hour, minute)
	}
}
