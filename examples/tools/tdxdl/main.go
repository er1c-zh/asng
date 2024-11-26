package main

import (
	"asng/models"
	"asng/proto"
	"context"
	"fmt"
	"os"
)

func main() {
	for i := 1; i < len(os.Args); i += 1 {

		var err error
		cli := proto.NewClient(context.Background(), proto.DefaultOption.
			WithDebugMode().
			WithTCPAddress("110.41.147.114:7709").
			WithDebugMode().
			WithMsgCallback(func(pi models.ProcessInfo) {
				fmt.Printf("%s\n", pi.Msg)
			}).WithMetaAddress("124.71.223.19:7727"))
		err = cli.Connect()
		if err != nil {
			fmt.Printf("error:%s", err)
			return
		}
		fmt.Printf("connected\n")
		data, err := cli.DownloadFile(os.Args[i])
		if err != nil {
			fmt.Printf("error:%s", err)
			return
		}
		err = os.WriteFile(os.Args[i], data, 0644)
		if err != nil {
			fmt.Printf("error:%s", err)
			return
		}
	}
}
