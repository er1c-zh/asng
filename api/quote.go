package api

import (
	"asng/models"
	"asng/proto"
	"fmt"
)

func (a *App) CandleStick(id models.StockIdentity, period proto.CandleStickPeriodType, cursor uint16) *proto.CandleStickResp {
	resp, err := a.cli.CandleStick(id, period, cursor)
	if err != nil {
		a.LogProcessError(models.ProcessInfo{Msg: fmt.Sprintf("candle stick failed: %s", err.Error())})
		return nil
	}
	return resp
}

func (a *App) TodayQuote(id models.StockIdentity) []proto.QuoteFrame {
	resp, err := a.cli.RealtimeGraph(id)
	if err != nil {
		a.LogProcessError(models.ProcessInfo{Msg: fmt.Sprintf("today quote failed: %s", err.Error())})
		return nil
	}
	return resp.List
}

func (a *App) RealtimeInfo(req []models.StockIdentity) *proto.RealtimeInfoResp {
	resp, err := a.cli.RealtimeInfo(req)
	if err != nil {
		return nil
	}
	return resp
}
