package proto

import (
	"asng/models"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func ParseBaseDBF(data []byte) (*models.BaseDBF, error) {
	DBF := &models.BaseDBF{}
	cursor := 0

	var err error

	cursor = 0x04
	DBF.RowCount, err = ReadInt(data, &cursor, DBF.RowCount)
	if err != nil {
		return nil, err
	}

	DBF.DataOffset, err = ReadInt(data, &cursor, DBF.DataOffset)
	if err != nil {
		return nil, err
	}
	DBF.ColMetaCount = (DBF.DataOffset - 1 - 0x20 /*header byte width*/) / 0x20 /* col meta byte width*/

	DBF.OneRowOffset, err = ReadInt(data, &cursor, DBF.OneRowOffset)
	if err != nil {
		return nil, err
	}

	cursor = 0x20
	for i := 0; i < int(DBF.BaseDBFHeader.ColMetaCount); i++ {
		c0 := cursor
		item := models.BaseDBFColMeta{}
		item.Name, err = ReadTDXString(data, &cursor, 0x0B)
		if err != nil {
			return nil, err
		}
		cursor += 0x01
		item.Offset, err = ReadInt(data, &cursor, item.Offset)
		if err != nil {
			return nil, err
		}
		cursor += 0x02
		item.Width, err = ReadInt(data, &cursor, item.Width)
		if err != nil {
			return nil, err
		}
		cursor = c0 + 0x20
		DBF.ColMetaList = append(DBF.ColMetaList, item)
	}
	{
		j, _ := json.MarshalIndent(DBF.ColMetaList, "", "  ")
		fmt.Printf("%s\n", j)
	}
	cursor = int(DBF.DataOffset)
	for i := 0; i < int(DBF.RowCount); i++ {
		item := models.BaseDBFRow{}
		item.Data = make(map[string]any)
		offset := int(DBF.DataOffset) + i*int(DBF.OneRowOffset)
		fmt.Printf("i: %d, offset: %d\n", i, offset)
		for j := 0; j < int(DBF.ColMetaCount); j++ {
			curData := data[offset+int(DBF.ColMetaList[j].Offset) : offset+int(DBF.ColMetaList[j].Offset+uint16(DBF.ColMetaList[j].Width))]
			curDataStr := strings.Trim(strings.TrimSpace(string(curData)), "\x00")
			switch DBF.ColMetaList[j].Name {
			case "GPDM":
				item.Code = string(data[offset+int(DBF.ColMetaList[j].Offset) : offset+int(DBF.ColMetaList[j].Offset+uint16(DBF.ColMetaList[j].Width))])
			case "SC":
				if err != nil {
					return nil, err
				}
				item.Market = models.MarketType(curDataStr[0] - '0')
			case "GXRQ", "SSDATE", "DY", "HY":
				if curDataStr == "" {
					continue
				}
				item.Data[DBF.ColMetaList[j].Name], err = strconv.ParseInt(curDataStr, 10, 64)
				if err != nil {
					return nil, err
				}
			default:
				if curDataStr == "" {
					continue
				}
				item.Data[DBF.ColMetaList[j].Name], err = strconv.ParseFloat(curDataStr, 64)
				if err != nil {
					fmt.Printf("error:%s\n", err)
					continue
				}
			}
		}

		DBF.Data = append(DBF.Data, item)
	}

	return DBF, nil
}
