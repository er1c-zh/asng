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
		if strings.HasPrefix(item.Code, s) {
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

func (a *App) StockMeta(s []string) map[string]models.StockMetaItem {
	result := make(map[string]models.StockMetaItem)
	for _, code := range s {
		if a.stockMetaMap[code] == nil {
			continue
		}
		result[code] = *a.stockMetaMap[code]
	}
	return result
}
