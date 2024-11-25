package api

import (
	"asng/models"
	"asng/proto"
	"fmt"
)

func (a *App) CandleStick(code string, period proto.CandleStickPeriodType, cursor uint16) *proto.CandleStickResp {
	meta, ok := a.stockMetaMap[code]
	if !ok {
		return nil
	}
	resp, err := a.cli.CandleStick(meta.Market, code, period, cursor)
	if err != nil {
		a.LogProcessError(models.ProcessInfo{Msg: fmt.Sprintf("candle stick failed: %s", err.Error())})
		return nil
	}
	return resp
}
