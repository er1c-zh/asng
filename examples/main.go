package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"asng/models"
	"asng/proto"
)

func main() {
	// testStockMeta()
	// testServerInfo()
	// testDownloadFile()
	// test0547()
	// testServerInfo()
	// test0537()
	testSubscribe()
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

	err = cli.Heartbeat()
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
	// Resp, err := cli.RealtimeInfo([]proto.StockQuery{{Market: uint8(models.MarketSH), Code: "999999"}})
	// If err != nil {
	// 	fmt.Printf("error:%s", err)
	// 	return
	// }
	// Fmt.Printf("%v\n", (resp.ItemList))
}

func testSubscribe() {
	var err error
	cli := proto.NewClient(context.Background(), proto.DefaultOption.
		WithDebugMode().
		WithTCPAddress("110.41.147.114:7709").
		WithDebugMode().
		WithMsgCallback(func(pi models.ProcessInfo) {
			fmt.Printf("%s\n", pi.Msg)
		}).WithMetaAddress("124.71.223.19:7727").
		WithRespHandler(proto.RealtimeSubscribeType0, proto.RealtimeSubscribeHandler))
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

	err = cli.Heartbeat()
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

	// err = cli.Nop(0x0547, 0x0029, "0100013630303034388ab00100") // subscribe
	// err = cli.Nop(0x0FC5, 0x1C00, "010036303030343800003200") // subscribe
	// err = cli.Nop(0x0FC5, 0x0100, "010036303030343800008403") // subscribe
	// err = cli.Nop(0x0FC6, 0x0101, "E9DA3401010036303030343800008403") // subscribe

	// resp, err := cli.RealtimeInfoSubscribe([]proto.StockQuery{
	// 	{Market: uint8(models.MarketSH), Code: "600000"},
	// 	{Market: uint8(models.MarketSZ), Code: "300050"},
	// })
	// resp, err := cli.RealtimeInfo([]proto.StockQuery{{Market: uint8(models.MarketSH), Code: "600000"}})
	// resp, err := cli.TXToday(models.StockIdentity{
	// 	MarketType: models.MarketSH, Code: "600000",
	// })
	// resp, err := cli.TXRealtime(models.StockIdentity{
	// MarketType: models.MarketSH, Code: "600000"})
	resp, err := cli.CandleStick(models.StockIdentity{MarketType: models.MarketSH, Code: "600000"}, proto.CandleStickPeriodType_Day, 0)
	if err != nil {
		fmt.Printf("error:%s", err)
		return
	}
	j, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Printf("%s\n", j)

	// resp, err := cli.RealtimeGraph(models.MarketSH, "600000")
	// if err != nil {
	// fmt.Printf("error:%s", err)
	// return
	// }

	time.Sleep(100 * time.Second)

}
