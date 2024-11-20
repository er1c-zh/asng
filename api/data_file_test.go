package api

import (
	"asng/models"
	"context"
	"encoding/json"
	"testing"
)

func TestFileManager(t *testing.T) {
	ctx := context.Background()
	fm, err := NewFileManager(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	err = SaveFile(fm, TypeStockMeta, models.StockMetaAll{})
	if err != nil {
		t.Error(err)
		return
	}
	h, d, err := ReadFile(fm, TypeStockMeta, models.StockMetaAll{})
	if err != nil {
		t.Error(err)
		return
	}
	j, _ := json.MarshalIndent(h, "", "  ")
	t.Log(string(j))
	j, _ = json.MarshalIndent(d, "", "  ")
	t.Log(string(j))
}
