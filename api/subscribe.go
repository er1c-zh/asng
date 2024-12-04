package api

import (
	"asng/models"
	"asng/proto"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type QuoteSubscripition struct {
	app    *App
	cli    *proto.Client
	m      map[string]*SubscribeReq
	ticker *time.Ticker
	urgent chan string
}

type SubscribeReq struct {
	Group     string
	Code      []string
	QuoteType string
}

func NewQuoteSubscripition(app *App) *QuoteSubscripition {
	return &QuoteSubscripition{
		app:    app,
		cli:    app.cli,
		m:      make(map[string]*SubscribeReq),
		ticker: time.NewTicker(time.Second * 3),
		urgent: make(chan string, 10),
	}
}

type QuoteSubscribeResp struct {
	RealtimeInfo proto.RealtimeInfoRespItem
	PriceFrame   models.QuoteFrameDataSingleValue
	VolumeFrame  models.QuoteFrameDataSingleValue
}

func (a *QuoteSubscripition) Start() {
	a.app.cli.RegisterRealtimeConsumer(func(d proto.RealtimeInfoRespItem) {
		runtime.EventsEmit(a.app.ctx, string(MsgKeySubscribeBroadcast),
			a.app.generateQuoteSubscribeResp(d))
	})
}

func (a *App) Subscribe(req []models.StockIdentity) []QuoteSubscribeResp {
	resp, err := a.cli.RealtimeInfoSubscribe(req)
	if err != nil {
		return nil
	}
	result := make([]QuoteSubscribeResp, 0)
	for _, v := range resp.ItemList {
		result = append(result, a.generateQuoteSubscribeResp(v))
	}
	return result
}

func (a *App) generateQuoteSubscribeResp(d proto.RealtimeInfoRespItem) QuoteSubscribeResp {
	return QuoteSubscribeResp{
		RealtimeInfo: d,
		PriceFrame: models.QuoteFrameDataSingleValue{
			QuoteFrame: models.QuoteFrame{
				TimeInMs: d.TickMilliSecTimestamp,
			},
			Value: d.CurrentPrice,
			Scale: int64(a.stockMetaMap[d.ID].Scale),
		},
		VolumeFrame: models.QuoteFrameDataSingleValue{
			QuoteFrame: models.QuoteFrame{
				TimeInMs: d.TickMilliSecTimestamp,
			},
			Value: d.CurrentVolume,
			Scale: 1,
		},
	}
}
