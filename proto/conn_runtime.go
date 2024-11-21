package proto

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/binary"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	maskRecvHeaderZipFlag = 0x10
)

type ConnRuntime struct {
	ctx context.Context

	seqID uint32

	heartbeatTicker *time.Ticker
	heartbeatFunc   func() error

	log func(format string, args ...any)

	done              chan struct{}
	muConn            sync.Mutex
	conn              net.Conn
	connected         bool
	sendCh            chan *reqPkg
	muHandlerRedister sync.Mutex
	handlerRegister   map[uint32]*reqPkg
}

type connRuntimeOpt struct {
	heartbeatInterval time.Duration
	heartbeatFunc     func() error
	log               func(format string, args ...any)
}

func newConnRuntime(ctx context.Context, opt connRuntimeOpt) *ConnRuntime {
	r := &ConnRuntime{
		log: func(format string, args ...any) {
			// do nothing
		},
	}

	r.ctx = ctx
	r.seqID = 0
	r.heartbeatTicker = time.NewTicker(opt.heartbeatInterval)
	r.heartbeatFunc = opt.heartbeatFunc
	if opt.log != nil {
		r.log = opt.log
	}
	r.done = make(chan struct{})
	r.connected = false
	r.sendCh = make(chan *reqPkg)
	r.handlerRegister = make(map[uint32]*reqPkg)

	go r.heartbeatTrigger()

	return r
}

func (r *ConnRuntime) heartbeatTrigger() {
	// TODO recover panic
	for {
		r.log("heartbeat ticker, seqID: %d", r.seqID)
		select {
		case <-r.heartbeatTicker.C:
			r.log("heartbeat %t", r.connected)
			if !r.connected {
				continue
			}
			t0 := time.Now()
			err := r.heartbeatFunc()
			if err != nil {
				r.resetConn()
				r.log("heartbeat fail: %v", err)
			} else {
				r.log("heartbeat success, cost: %d ms", time.Since(t0).Milliseconds())
			}
		case <-r.done:
			return
		}
	}
}

func (r *ConnRuntime) connect(addr string) error {
	r.muConn.Lock()
	defer r.muConn.Unlock()
	var err error

	r.conn, err = net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go r.sendHandler()
	go r.recvHandler()

	err = r.heartbeatFunc()
	if err != nil {
		return err
	}
	r.connected = true
	r.log("connect to %s success", addr)
	return nil
}

func (r *ConnRuntime) sendHandler() {
	r.log("send handler start")
	for {
		d := <-r.sendCh
		if d == nil {
			continue
		}
		r.log("send: %s %d", d.header.Method, d.header.SeqID)
		r.muHandlerRedister.Lock()
		r.handlerRegister[uint32(d.header.SeqID)<<16+uint32(d.header.Type0)] = d
		r.muHandlerRedister.Unlock()
		n, err := r.conn.Write(d.body)
		if err != nil {
			r.log("write fail: %s", err.Error())
			break
		}
		r.log("send: %d, size: %d", d.header.SeqID, n)
	}
	r.resetConn()
}

func (r *ConnRuntime) recvHandler() {
	r.log("read start.")
	for {
		var err error

		// read header
		headerBuf := make([]byte, 16)
		_, err = io.ReadFull(r.conn, headerBuf)
		if err != nil {
			r.log("read header fail: %v", err)
			break
		}
		var header RespHeader
		if err := binary.Read(bytes.NewBuffer(headerBuf), binary.LittleEndian, &header); err != nil {
			r.log("parse header fail: %v", err)
			break
		}

		r.log("read: %d, size: %d", header.SeqID, header.PkgDataSize)

		body := make([]byte, header.PkgDataSize)
		_, err = io.ReadFull(r.conn, body)
		if err != nil {
			r.log("read body fail: %v", err)
			break
		}
		if header.Flag&maskRecvHeaderZipFlag > 0 {
			zipReader, err := zlib.NewReader(bytes.NewReader(body))
			if err != nil {
				r.log("unzip fail: %v", err)
				continue
			}
			body, err = io.ReadAll(zipReader)
			if err != nil {
				r.log("unzip fail: %v", err)
				_ = zipReader.Close()
				continue
			}
			_ = zipReader.Close()
		}
		// dispatch
		r.muHandlerRedister.Lock()
		reqPkg, ok := r.handlerRegister[uint32(header.SeqID)<<16+uint32(header.Type0)]
		if !ok {
			r.muHandlerRedister.Unlock()
			r.log("handler not found: %s %d", header.Method, header.SeqID)
			continue
		}
		delete(r.handlerRegister, uint32(header.SeqID)<<16+uint32(header.Type0))
		r.muHandlerRedister.Unlock()
		if reqPkg == nil || reqPkg.callback == nil {
			r.log("callback is nil: %s %d", header.Method, header.SeqID)
			continue
		}
		r.log("call callback: %s %d", header.Method, header.SeqID)
		reqPkg.callback <- &respPkg{header: header, body: body}
	}
	r.resetConn()
}

func (r *ConnRuntime) resetConn() {
	r.log("reset conn")
	r.muConn.Lock()
	defer r.muConn.Unlock()
	if r.conn != nil {
		_ = r.conn.Close()
	}
	r.connected = false
}

func (r *ConnRuntime) genSeqID() uint16 {
	return uint16(atomic.AddUint32(&r.seqID, 1))
}

func (r *ConnRuntime) resetSeqID(base uint16) {
	atomic.StoreUint32(&r.seqID, uint32(base))
}

func (r *ConnRuntime) isConnected() bool {
	return r.connected
}
