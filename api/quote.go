package api

import (
	"asng/models"
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
	Price    []models.QuoteFrameDataSingleValue
	AvgPrice []models.QuoteFrameDataSingleValue
	Volume   []models.QuoteFrameDataSingleValue
}

func (a *App) TodayQuote(id models.StockIdentity) TodayQuoteResp {
	resp := TodayQuoteResp{}
	var err error
	tx, err := a.cli.TXToday(id)
	if err != nil {
		a.LogProcessError(models.ProcessInfo{Msg: fmt.Sprintf("today quote failed: %s", err.Error())})
		resp.Code = 1
		resp.Message = "fail"
		return resp
	}
	resp.Price = make([]models.QuoteFrameDataSingleValue, 0)
	resp.AvgPrice = make([]models.QuoteFrameDataSingleValue, 0)
	resp.Volume = make([]models.QuoteFrameDataSingleValue, 0)

	t0 := utils.GetTodayWithOffset(0, 0, 0)
	txOffsetInOneMinute := 0
	tickCur := uint16(0)

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
		resp.Price = append(resp.Price, models.QuoteFrameDataSingleValue{
			QuoteFrame: f.Clone().SetType(models.QuoteTypeLine),
			Value:      item.Price,
			Scale:      10000,
		})
		// resp.AvgPrice = append(resp.AvgPrice, models.QuoteFrameDataSingleValue{
		// 	QuoteFrame: f.Clone().SetType(models.QuoteTypeLine),
		// 	Value:      item.AvgPrice,
		// })
		// resp.Volume = append(resp.Volume, models.QuoteFrameDataSingleValue{
		// 	QuoteFrame: f.Clone().SetType(models.QuoteTypeBar),
		// 	Value:      item.Volume,
		// })
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
