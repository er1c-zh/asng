package proto

import (
	"asng/models"
	"time"
)

const (
	_defaultTCPAddress        = "119.147.212.81:7709"
	_defaultRetryTimes        = 3
	_defaultHeartbeatInterval = 15 * time.Second
	_defaultMetaAddress       = "124.71.223.19:7727"
)

type Options struct {
	TCPAddress        string   // 服务器地址
	TCPAddressPool    []string // 服务器地址池
	MaxRetryTimes     int      // 重试次数
	HeartbeatInterval time.Duration
	MsgCallback       func(models.ProcessInfo)
	Debug             bool
	MetaAddress       string // meta 数据服务器
	RespHandler       map[uint16]RespHandler
}

func defaultOptions() *Options {
	return &Options{
		TCPAddress:        _defaultTCPAddress,
		MaxRetryTimes:     _defaultRetryTimes,
		HeartbeatInterval: _defaultHeartbeatInterval,
		MsgCallback: func(pi models.ProcessInfo) {
			// do nothing
		},
		Debug:       false,
		MetaAddress: _defaultMetaAddress,
	}
}

func ApplyOptions(opts ...Option) *Options {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	return o
}

type Option func(options *Options)

var DefaultOption Option = func(options *Options) {}

func (Option Option) WithTCPAddress(tcpAddress string) Option {
	return func(o *Options) {
		o.TCPAddress = tcpAddress
		Option(o)
	}
}

func (Option Option) WithTCPAddressPool(ips ...string) Option {
	return func(o *Options) {
		o.TCPAddressPool = ips
		Option(o)
	}
}

func (Option Option) WithHeartbitInterval(interval time.Duration) Option {
	return func(o *Options) {
		o.HeartbeatInterval = interval
		Option(o)
	}
}

func (Option Option) WithMsgCallback(callback func(models.ProcessInfo)) Option {
	return func(o *Options) {
		o.MsgCallback = callback
		Option(o)
	}
}

func (Option Option) WithDebugMode() Option {
	return func(o *Options) {
		o.Debug = true
		Option(o)
	}
}

func (Option Option) WithMetaAddress(metaAddress string) Option {
	return func(o *Options) {
		o.MetaAddress = metaAddress
		Option(o)
	}
}

func (Option Option) WithRespHandler(type0 uint16, h RespHandler) Option {
	return func(o *Options) {
		if o.RespHandler == nil {
			o.RespHandler = make(map[uint16]RespHandler, 0)
		}
		o.RespHandler[type0] = h
		Option(o)
	}
}
