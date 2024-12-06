package api

import (
	"asng/models"
	"asng/models/value"
	"asng/proto"
	"asng/utils"
	"fmt"
	"time"
)

func (a *App) CandleStick(id models.StockIdentity, period proto.CandleStickPeriodType, cursor uint16) *proto.CandleStickResp {
	resp, err := a.cli.CandleStick(id, period, cursor)
	if err != nil {
		a.LogProcessError(models.ProcessInfo{Msg: fmt.Sprintf("candle stick failed: %s", err.Error())})
		return nil
	}
	return resp
}

type TodayQuoteResp struct {
	models.BaseResp
	Frames []models.QuoteFrameRealtime
}

func (a *App) TodayQuote(id models.StockIdentity) TodayQuoteResp {
	resp := TodayQuoteResp{}
	resp.BaseResp = models.BaseResp{
		Code:    0,
		Message: "success",
	}

	var err error
	tx, err := a.cli.TXToday(id)
	if err != nil {
		a.LogProcessError(models.ProcessInfo{Msg: fmt.Sprintf("today quote failed: %s", err.Error())})
		resp.Code = 1
		resp.Message = "fail"
		return resp
	}

	resp.Frames = make([]models.QuoteFrameRealtime, 0, len(tx))

	t0 := utils.GetTodayWithOffset(0, 0, 0)
	txOffsetInOneMinute := 0
	tickCur := uint16(0)

	var (
		totalVolume int64 = 0
		totalSum    int64 = 0
	)
	for _, item := range tx {
		if item.Tick == tickCur {
			txOffsetInOneMinute += 3 // 3 seconds per tick
		} else {
			tickCur = item.Tick
			txOffsetInOneMinute = 0
		}
		f := models.QuoteFrame{
			TimeInMs: t0.Add(
				time.Duration(item.Tick)*time.Minute +
					time.Duration(txOffsetInOneMinute)*time.Second).
				UnixMilli(),
			TimeSpanInMs: time.Second.Milliseconds(),
		}

		totalVolume += item.Volume
		totalSum += item.Price * item.Volume

		resp.Frames = append(resp.Frames, models.QuoteFrameRealtime{
			QuoteFrame: f,
			Price:      value.IntWithScale(item.Price, 10000),
			AvgPrice:   value.IntWithScale(totalSum/totalVolume, 10000),
			Volume:     value.IntWithScale(item.Volume, 1),
		})
	}
	return resp
}

func (a *App) RealtimeInfo(req []models.StockIdentity) *proto.RealtimeInfoResp {
	resp, err := a.cli.RealtimeInfo(req)
	if err != nil {
		return nil
	}
	return resp
}
