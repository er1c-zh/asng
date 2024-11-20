package api

import (
	"asng/models"
	"asng/proto"
)

type MsgKey string

const (
	MsgKeyInit               MsgKey = "init"
	MsgKeyProcessMsg         MsgKey = "processMsg"
	MsgKeyServerStatus       MsgKey = "serverStatus"
	MsgKeySubscribeBroadcast MsgKey = "subscribeBroadcast"
)

var ExportMsg = []struct {
	Value  MsgKey
	TSName string
}{
	{MsgKeyInit, string(MsgKeyInit)},
	{MsgKeyProcessMsg, string(MsgKeyProcessMsg)},
	{MsgKeyServerStatus, string(MsgKeyServerStatus)},
	{MsgKeySubscribeBroadcast, string(MsgKeySubscribeBroadcast)},
}

var ExportMarketType = []struct {
	Value  models.MarketType
	TSName string
}{
	{models.MarketSZ, models.MarketSZ.String()},
	{models.MarketSH, models.MarketSH.String()},
	{models.MarketBJ, models.MarketBJ.String()},
}

var ExportCandleStickPeriodType = []struct {
	Value  proto.CandleStickPeriodType
	TSName string
}{
	{proto.CandleStickPeriodType_1Min, "CandleStickPeriodType1Min"},
	{proto.CandleStickPeriodType_5Min, "CandleStickPeriodType5Min"},
	{proto.CandleStickPeriodType_15Min, "CandleStickPeriodType15Min"},
	{proto.CandleStickPeriodType_30Min, "CandleStickPeriodType30Min"},
	{proto.CandleStickPeriodType_1Hour, "CandleStickPeriodType60Min"},
	{proto.CandleStickPeriodType_Day, "CandleStickPeriodType1Day"},
	{proto.CandleStickPeriodType_Week, "CandleStickPeriodType1Week"},
	{proto.CandleStickPeriodType_Month, "CandleStickPeriodType1Month"},
}
