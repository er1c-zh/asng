package proto

import (
	"asng/models"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand/v2"
	"time"
)

type Client struct {
	ctx context.Context
	opt *Options

	codec *tdxCodec

	dataConn *ConnRuntime

	done chan struct{}

	// TDX seed
	handShakeSeed uint32
	macAddr       [6]byte

	persistenceHandlerChan chan *respPkg
	persistenceHandler     map[uint16]RespHandler
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

type RespHandler func(ctx context.Context, h *RespHeader, data []byte)

// private
func (c *Client) init() error {
	var err error
	c.done = make(chan struct{})
	c.codec, err = NewTDXCodec()
	if err != nil {
		return err
	}
	c.handShakeSeed = rand.Uint32()%30000 + 1
	c.Log("handshake seed: %d", c.handShakeSeed)
	c.macAddr, err = GenerateMACAddr()
	if err != nil {
		return err
	}
	c.Log("mac addr: %s", hex.EncodeToString(c.macAddr[:]))
	c.persistenceHandlerChan = make(chan *respPkg, 1024)
	c.persistenceHandler = make(map[uint16]RespHandler)
	for k, v := range c.opt.RespHandler {
		c.persistenceHandler[k] = v
	}
	go c.eventloop()
	return nil
}

func (c *Client) eventloop() {
	for {
		select {
		case <-c.done:
			return
		case pkg := <-c.persistenceHandlerChan:
			if h, ok := c.opt.RespHandler[pkg.header.Type0]; ok {
				go func() {
					h(c.ctx, &pkg.header, pkg.body)
				}()
			}
		}
	}
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

	t0 := time.Now()
	var err error
	reqHeader := ReqHeader{
		MagicNumber: 0x0C,
		HeaderIdentifer: HeaderIdentifer{
			SeqID: 0,
			Group: 0,
			Type0: 0,
		},
		PacketType: 0x01,
		PkgLen1:    0,
		PkgLen2:    0,
		Method:     0,
	}

	defer func() {
		c.LogDebug("do %s with err: %v, cost: %d ms", reqHeader.Method.String(), err, time.Since(t0).Milliseconds())
	}()

	err = api.FillReqHeader(c.ctx, &reqHeader)
	if err != nil {
		return err
	}

	reqHeader.SeqID = conn.genSeqID(reqHeader.Group)

	reqData, err := api.MarshalReqBody(c.ctx)
	if err != nil {
		return err
	}

	if api.NeedEncrypt(c.ctx) {
		reqData, err = c.codec.Encode(reqData)
		if err != nil {
			return err
		}
	}

	reqHeader.PkgLen1 = 2 + uint16(len(reqData))
	reqHeader.PkgLen2 = 2 + uint16(len(reqData))

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
		c.LogDebug("send %s\n%s", reqHeader.Method, hex.Dump(reqBuf.Bytes()))
	}

	callback := make(chan *respPkg)
	conn.sendCh <- &reqPkg{header: reqHeader, body: reqBuf.Bytes(), callback: callback}
	// wait resp
	respPkg := <-callback

	if c.opt.Debug && api.IsDebug(c.ctx) {
		c.LogDebug("recv %s\n%s", respPkg.header.Method, hex.Dump(respPkg.body))
	}

	err = api.UnmarshalResp(c.ctx, respPkg.body)
	if err != nil {
		return err
	}
	return nil
}

// public
func (c *Client) Connect() error {
	var err error
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
	err = c.dataConn.connect(c.opt.TCPAddress)
	if err != nil {
		return err
	}

	for type0 := range c.persistenceHandler {
		c.dataConn.RegisterHandler(type0, c.persistenceHandlerChan)
	}

	_, err = c.TDXHandshake()
	if err != nil {
		defer c.Disconnect()
		return err
	}
	return nil
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
	close(c.done)
	if c.dataConn != nil && c.dataConn.isConnected() {
		c.dataConn.resetConn()
	}
	return nil
}
