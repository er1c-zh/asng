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
	test052D()
}

func test052D() {
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

	resp, err := cli.CandleStick(models.MarketSH, "600000", proto.CandleStickPeriodType_Day, 0)
	if err != nil {
		fmt.Printf("error:%s", err)
		return
	}
	j, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Printf("%s\n", j)
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
