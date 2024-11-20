package proto

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
)

const (
	DownloadFileBatchSize = 30000
)

func (c *Client) DownloadFile(fileName string) ([]byte, error) {
	var err error
	fileBuf := bytes.NewBuffer(nil)
	offset := 0
	d := DownloadFile{}
	// d.SetDebug(c.ctx)
	d.Req.FileName = [300]byte{}
	copy(d.Req.FileName[:], fileName)
	d.Req.Size = DownloadFileBatchSize
	for {
		d.Req.Offset = uint32(offset)
		err = do(c, c.dataConn, &d)
		if err != nil {
			return nil, err
		}
		fileBuf.Write(d.Resp.Data)
		if d.Resp.Size < DownloadFileBatchSize {
			break
		} else {
			offset += int(d.Resp.Size)
		}
	}
	return fileBuf.Bytes(), nil
}

type DownloadFile struct {
	BlankCodec
	Req  DownloadFileReq
	Resp DownloadFileResp
}
type DownloadFileReq struct {
	Offset   uint32
	Size     uint32
	FileName [300]byte
}
type DownloadFileResp struct {
	Size uint32
	Data []byte
}

func (obj *DownloadFile) FillReqHeader(ctx context.Context, header *ReqHeader) error {
	header.MagicNumber = 0
	header.SeqID = 0
	header.PacketType = 0
	header.Method = 0x06B9
	header.PacketType = 2
	return nil
}

func (obj *DownloadFile) MarshalReqBody(ctx context.Context) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := binary.Write(buf, binary.LittleEndian, obj.Req)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (obj *DownloadFile) UnmarshalResp(ctx context.Context, data []byte) error {
	var err error
	cursor := 0
	obj.Resp.Size, err = ReadInt(data, &cursor, obj.Resp.Size)
	if err != nil {
		return err
	}

	// consider performence copy directly
	if len(data) < cursor+int(obj.Resp.Size) {
		return fmt.Errorf("data not enough")
	}
	obj.Resp.Data = make([]byte, obj.Resp.Size)
	copy(obj.Resp.Data, data[cursor:cursor+int(obj.Resp.Size)])

	return nil
}
