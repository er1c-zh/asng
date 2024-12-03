package api

import (
	"asng/models"
	"strings"
)

func (c *App) CommandMatch(s string) []models.StockMetaItem {
	result := []models.StockMetaItem{}
	if c.stockMeta == nil {
		return result
	}
	for _, item := range c.stockMeta.StockList {
		if strings.HasPrefix(item.ID.Code, s) {
			result = append(result, item)
		} else if strings.HasPrefix(item.PinYinInitial, s) {
			result = append(result, item)
		}
		if len(result) > 5 {
			break
		}
	}
	return result
}

func (a *App) ServerStatus() *models.ServerStatus {
	return a.status
}

func (a *App) StockMeta(s []models.StockIdentity) []models.StockMetaItem {
	result := make([]models.StockMetaItem, 0)
	for _, id := range s {
		if a.stockMetaMap[id] == nil {
			continue
		}
		result = append(result, *a.stockMetaMap[id])
	}
	return result
}
