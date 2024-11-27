package api

import (
	"asng/proto"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type QuoteSubscripition struct {
	app    *App
	rwMu   sync.RWMutex
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

func (a *QuoteSubscripition) Start() {
	a.app.cli.RegisterRealtimeConsumer(func(d proto.RealtimeInfoRespItem) {
		runtime.EventsEmit(a.app.ctx, string(MsgKeySubscribeBroadcast), d)
	})
}

func (a *App) Subscribe(req []proto.StockQuery) *proto.RealtimeInfoResp {
	resp, err := a.cli.RealtimeInfoSubscribe(req)
	if err != nil {
		return nil
	}
	return resp
}
