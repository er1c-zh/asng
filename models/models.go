package models

type ProcessInfoType uint8

const (
	ProcessInfoTypeDebug ProcessInfoType = 0
	ProcessInfoTypeInfo  ProcessInfoType = 1
	ProcessInfoTypeWarn  ProcessInfoType = 2
	ProcessInfoTypeErr   ProcessInfoType = 3
)

type ProcessInfo struct {
	Type ProcessInfoType
	Msg  string
}

type ServerStatus struct {
	Connected  bool
	ServerInfo string
}

type MarketType uint16

const (
	MarketSZ MarketType = 0
	MarketSH MarketType = 1
	MarketBJ MarketType = 2
)

func (mt MarketType) Uint8() uint8 {
	return uint8(mt)
}

func (m MarketType) String() string {
	switch m {
	case MarketSZ:
		return "深圳"
	case MarketSH:
		return "上海"
	case MarketBJ:
		return "北京"
	default:
		return "unknown"
	}
}

type StockMetaAll struct {
	StockList []StockMetaItem
}

type StockMetaItem struct {
	ID             StockIdentity
	Desc           string
	PinYinInitial  string // TODO multi-prounciation
	Scale          uint16
	F0             float64
	YesterdayClose float64
	BlockID        uint16
	F3             uint16

	BaseDBFItem *BaseDBFRow
}

type BaseDBF struct {
	BaseDBFHeader
	Data []BaseDBFRow
}

type BaseDBFHeader struct {
	RowCount     uint32
	DataOffset   uint16
	OneRowOffset uint16
	ColMetaCount uint16
	ColMetaList  []BaseDBFColMeta
}

type BaseDBFColMeta struct {
	Name   string
	Offset uint16
	Width  uint8
}
type BaseDBFRow struct {
	ID   StockIdentity
	Data map[string]any
}

type BaseReq struct{}

type BaseResp struct {
	Code    int
	Message string
}
