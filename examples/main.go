package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"asng/models"
	"asng/proto"
)

func main() {
	// testStockMeta()
	// testServerInfo()
	// testDownloadFile()
	// test0547()
	// testServerInfo()
	test000D()
}

func test000D() {
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

	_, err = cli.TDXHandshake()
	if err != nil {
		fmt.Printf("error:%s", err)
		return
	}

	serverInfo, err := cli.ServerInfo()
	if err != nil {
		fmt.Printf("error:%s", err)
		return
	}
	fmt.Printf("server name:%s\n", serverInfo.Name)

	// err = cli.Nop(0x0fdb, 0x0099, "7464786c6576656c000000b81ef540070000000000000000000000000005")
	// if err != nil {
	// 	fmt.Printf("error:%s", err)
	// 	return
	// }

	// err = cli.Nop(0x0002, 0x0095, "")
	// if err != nil {
	// 	fmt.Printf("error:%s", err)
	// 	return
	// }

	// err = cli.Nop(0x000a, 0x0093, "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	// if err != nil {
	// 	fmt.Printf("error:%s", err)
	// 	return
	// }

	// err = cli.Nop(0x0fde, 0x0071, "")
	// if err != nil {
	// 	fmt.Printf("error:%s", err)
	// 	return
	// }

	// cli.ResetDataConnSeqID(0x0000)
	resp, err := cli.Realtime([]proto.StockQuery{{Market: uint8(models.MarketSH), Code: "999999"}})
	if err != nil {
		fmt.Printf("error:%s", err)
		return
	}
	fmt.Printf("%v\n", (resp.ItemList))
}

func test0547() {
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

	// cli.Heartbeat()

	resp, err := cli.Realtime([]proto.StockQuery{})
	if err != nil {
		fmt.Printf("error:%s", err)
		return
	}
	fmt.Printf("%s\n", hex.Dump(resp.Data))
	j, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Printf("%s\n", j)
}
