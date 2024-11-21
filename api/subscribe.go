package api

import (
	"asng/models"
	"asng/proto"
	"fmt"
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

func (a *QuoteSubscripition) Subscribe(req *SubscribeReq) {
	a.rwMu.Lock()
	defer a.rwMu.Unlock()
	old, ok := a.m[req.Group]
	if ok {
		old.Code = append(old.Code, req.Code...)
	} else {
		a.m[req.Group] = req
	}
	a.urgent <- req.Group
}

func (a *QuoteSubscripition) Unsubscribe(req *SubscribeReq) {
	a.rwMu.Lock()
	defer a.rwMu.Unlock()
	old, ok := a.m[req.Group]
	if !ok {
		return
	}
	if len(req.Code) == 0 {
		delete(a.m, req.Group)
		return
	}
	b := make([]string, len(old.Code)-1)
	for _, code := range old.Code {
		for _, c := range req.Code {
			if code == c {
				goto next
			}
		}
		b = append(b, code)
	next:
	}
	old.Code = b
}

func (a *QuoteSubscripition) Start() {
	go a.worker()
}

func (a *QuoteSubscripition) worker() {
	for {
		select {
		case <-a.ticker.C:
			a.rwMu.RLock()
			for _, req := range a.m {
				realtimeReq := make([]proto.StockQuery, 0, len(req.Code))
				for _, code := range req.Code {
					if a.app.stockMetaMap[code] == nil {
						continue
					}
					realtimeReq = append(realtimeReq,
						proto.StockQuery{Market: uint8(a.app.stockMetaMap[code].Market), Code: code})
				}
				go func(req []proto.StockQuery, group string) {
					resp, err := a.cli.RealtimeInfo(req)
					if err != nil {
						a.app.LogProcessError(models.ProcessInfo{Msg: fmt.Sprintf("realtime subscribe failed: %s", err.Error())})
						return
					}
					runtime.EventsEmit(a.app.ctx, string(MsgKeySubscribeBroadcast), group, resp.ItemList)
				}(realtimeReq, req.Group)
			}
			a.rwMu.RUnlock()
		case group := <-a.urgent:
			a.rwMu.RLock()
			req, ok := a.m[group]
			if ok && req != nil {
				realtimeReq := make([]proto.StockQuery, 0, len(req.Code))
				for _, code := range req.Code {
					if a.app.stockMetaMap[code] == nil {
						continue
					}
					realtimeReq = append(realtimeReq,
						proto.StockQuery{Market: uint8(a.app.stockMetaMap[code].Market), Code: code})
				}
				go func(req []proto.StockQuery, group string) {
					resp, err := a.cli.RealtimeInfo(req)
					if err != nil {
						a.app.LogProcessError(models.ProcessInfo{Msg: fmt.Sprintf("realtime subscribe failed: %s", err.Error())})
						return
					}
					runtime.EventsEmit(a.app.ctx, string(MsgKeySubscribeBroadcast), group, resp.ItemList)
				}(realtimeReq, req.Group)
			}
			a.rwMu.RUnlock()
		}
	}
}

func (a *App) Subscribe(req *SubscribeReq) {
	a.qs.Subscribe(req)
}

func (a *App) Unsubscribe(req *SubscribeReq) {
	a.qs.Unsubscribe(req)
}
