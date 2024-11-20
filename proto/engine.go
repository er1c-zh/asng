package proto

import (
	"asng/models"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"
)

type Client struct {
	ctx context.Context
	opt *Options

	codec *tdxCodec

	dataConn *ConnRuntime
}

type reqPkg struct {
	body     []byte
	header   ReqHeader
	callback chan *respPkg
}
type respPkg struct {
	header RespHeader
	body   []byte
}

func NewClient(ctx context.Context, opt Option) *Client {
	cli := &Client{
		opt: ApplyOptions(opt),
		ctx: ctx,
	}
	cli.init()
	cli.Log("init success.")
	return cli
}

// private
func (c *Client) init() error {
	var err error
	c.codec, err = NewTDXCodec()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Log(msg string, args ...any) {
	c.opt.MsgCallback(models.ProcessInfo{Type: models.ProcessInfoTypeInfo, Msg: fmt.Sprintf(msg, args...)})
}

func (c *Client) LogDebug(msg string, args ...any) {
	c.opt.MsgCallback(models.ProcessInfo{Type: models.ProcessInfoTypeDebug, Msg: fmt.Sprintf(msg, args...)})
}

// use generic type
func do[T Codec](c *Client, conn *ConnRuntime, api T) error {
	if c == nil {
		return fmt.Errorf("client is nil")
	}
	if conn == nil {
		return fmt.Errorf("conn is nil or not connected")
	}
	c.LogDebug("start do")
	var err error
	reqHeader := ReqHeader{
		MagicNumber: 0x0C,
		SeqID:       conn.genSeqID(),
		PacketType:  0x01,
		PkgLen1:     0,
		PkgLen2:     0,
		Method:      0,
	}
	err = api.FillReqHeader(c.ctx, &reqHeader)
	if err != nil {
		return err
	}
	reqData, err := api.MarshalReqBody(c.ctx)
	if err != nil {
		return err
	}
	reqHeader.PkgLen1 = 2 + uint16(len(reqData))
	reqHeader.PkgLen2 = 2 + uint16(len(reqData))

	if api.NeedEncrypt(c.ctx) {
		reqData, err = c.codec.Encode(reqData)
		if err != nil {
			return err
		}
	}

	// send req
	reqBuf := bytes.NewBuffer(nil)
	err = binary.Write(reqBuf, binary.LittleEndian, reqHeader)
	if err != nil {
		return err
	}
	_, err = reqBuf.Write(reqData)
	if err != nil {
		return err
	}

	if c.opt.Debug && api.IsDebug(c.ctx) {
		c.LogDebug("send %s", reqHeader.Method)
		c.LogDebug("%s", hex.Dump(reqBuf.Bytes()))
	}

	callback := make(chan *respPkg)
	conn.sendCh <- &reqPkg{header: reqHeader, body: reqBuf.Bytes(), callback: callback}
	// wait resp
	respPkg := <-callback

	if c.opt.Debug && api.IsDebug(c.ctx) {
		c.LogDebug("recv %s: %s", respPkg.header.Method, hex.EncodeToString(respPkg.body))
	}

	err = api.UnmarshalResp(c.ctx, respPkg.body)
	if err != nil {
		return err
	}
	return nil
}

// public
func (c *Client) Connect() error {
	if c.dataConn != nil && c.dataConn.isConnected() {
		return nil
	}
	if c.dataConn == nil {
		c.dataConn = newConnRuntime(c.ctx, connRuntimeOpt{
			heartbeatInterval: c.opt.HeartbeatInterval,
			log:               c.LogDebug,
			heartbeatFunc: func() error {
				return c.Heartbeat()
			},
		})
	}
	return c.dataConn.connect(c.opt.TCPAddress)
}

func (c *Client) NewMetaConnection() (*ConnRuntime, error) {
	conn := newConnRuntime(c.ctx, connRuntimeOpt{
		heartbeatInterval: c.opt.HeartbeatInterval,
		log:               c.Log,
		heartbeatFunc:     func() error { return nil },
	})
	err := conn.connect(c.opt.MetaAddress)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *Client) Disconnect() error {
	panic("implement me!")
}

func (c *Client) Handshake() error {
	handshake := &Handshake{}
	err := do(c, c.dataConn, handshake)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Heartbeat() error {
	t0 := time.Now()
	heartbeat := &Heartbeat{}
	err := do(c, c.dataConn, heartbeat)
	if err != nil {
		return err
	}
	c.Log("heartbeat success, cost: %d ms", time.Since(t0).Milliseconds())
	return nil
}
