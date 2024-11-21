package api

import (
	"asng/models"
	"asng/proto"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
	fm  *FileManager

	initOnce sync.Once
	initDone bool

	status *models.ServerStatus

	// data in memory
	stockMeta    *models.StockMetaAll
	stockMetaMap map[string]*models.StockMetaItem // WARN: read-only after init

	cli *proto.Client

	// quote subscription
	qs *QuoteSubscripition
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.status = &models.ServerStatus{}
}

func (a *App) Shutdown(_ctx context.Context) {
	if a.cli != nil {
		a.cli.Disconnect()
	}
}

func (a *App) EmitProcessInfo(i models.ProcessInfo) {
	runtime.EventsEmit(a.ctx, string(MsgKeyProcessMsg), i)
}

func (a *App) LogProcessInfo(i models.ProcessInfo) {
	i.Type = models.ProcessInfoTypeInfo
	a.EmitProcessInfo(i)
}

func (a *App) LogProcessWarn(i models.ProcessInfo) {
	i.Type = models.ProcessInfoTypeWarn
	a.EmitProcessInfo(i)
}

func (a *App) LogProcessError(i models.ProcessInfo) {
	i.Type = models.ProcessInfoTypeErr
	a.EmitProcessInfo(i)
}

const (
	InitDone  = "done"
	InitStart = "start"
)

func (a *App) Init() {
	go a.asyncInit()
}

func (a *App) asyncInit() {
	a.initOnce.Do(func() {
		var err error

		a.LogProcessInfo(models.ProcessInfo{Msg: "initializing..."})

		a.fm, err = NewFileManager(a.ctx)
		if err != nil {
			a.LogProcessError(models.ProcessInfo{Msg: fmt.Sprintf("file manager failed: %s", err.Error())})
			return
		}

		{
			a.LogProcessInfo(models.ProcessInfo{Msg: "initializing client..."})
			a.cli = proto.NewClient(a.ctx, proto.DefaultOption.
				WithDebugMode().
				WithTCPAddress("110.41.147.114:7709").
				WithMsgCallback(a.EmitProcessInfo).
				WithMetaAddress("124.71.223.19:7727"))
			err = a.cli.Connect()
			if err != nil {
				a.LogProcessError(models.ProcessInfo{Msg: fmt.Sprintf("connect client failed: %s", err.Error())})
				return
			}
			serverInfo, err := a.cli.ServerInfo()
			if err != nil {
				a.LogProcessError(models.ProcessInfo{Msg: fmt.Sprintf("get server info failed: %s", err.Error())})
				return
			}
			a.updateServerStatus(func(ss *models.ServerStatus) {
				ss.Connected = true
				ss.ServerInfo = serverInfo.Name
			})
		}

		{
			a.LogProcessInfo(models.ProcessInfo{Msg: "initializing file manager..."})
			a.fm, err = NewFileManager(a.ctx)
			if err != nil {
				a.LogProcessError(models.ProcessInfo{Msg: fmt.Sprintf("file manager failed: %s", err.Error())})
				return
			}
		}

		{
			a.LogProcessInfo(models.ProcessInfo{Msg: "loading stock meta..."})
			t0 := time.Now()
			_, a.stockMeta, err = a.fm.LoadStockMeta()
			if err != nil {
				a.LogProcessWarn(models.ProcessInfo{Msg: fmt.Sprintf("read stock meta failed: %s", err.Error())})
			}
			if a.stockMeta == nil {
				a.LogProcessInfo(models.ProcessInfo{Msg: "stock meta not found, loading from server..."})
				a.stockMeta, err = a.cli.StockMetaAll()
				if err != nil {
					a.LogProcessError(models.ProcessInfo{Msg: fmt.Sprintf("read stock meta from server failed: %s", err.Error())})
					return
				}
				a.LogProcessInfo(models.ProcessInfo{Msg: "stock meta saving..."})
				err = a.fm.SaveStockMeta(a.stockMeta)
				if err != nil {
					a.LogProcessError(models.ProcessInfo{Msg: fmt.Sprintf("save stock meta failed: %s", err.Error())})
					return
				}
			}
			a.stockMetaMap = make(map[string]*models.StockMetaItem, len(a.stockMeta.StockList))
			for _, v := range a.stockMeta.StockList {
				a.stockMetaMap[v.Code] = &v
			}
			a.LogProcessInfo(models.ProcessInfo{Msg: fmt.Sprintf("load stock meta cost: %d ms", time.Since(t0).Milliseconds())})
		}
		{
			a.LogProcessInfo(models.ProcessInfo{Msg: "initializing quote subscription..."})
			a.qs = NewQuoteSubscripition(a)
			a.qs.Start()
			a.LogProcessInfo(models.ProcessInfo{Msg: "quote subscription initialized"})
		}
		a.initDone = true
	})
	runtime.EventsEmit(a.ctx, string(MsgKeyInit), a.initDone)
}

func (a *App) updateServerStatus(f func(*models.ServerStatus)) {
	f(a.status)
	runtime.EventsEmit(a.ctx, string(MsgKeyServerStatus), a.status)
}

type ExportStruct struct {
	F0 models.ServerStatus
	F1 []proto.RealtimeRespItem
}

func (a *App) MakeWailsHappy() ExportStruct {
	return ExportStruct{}
}
