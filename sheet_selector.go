package excelstreamer

import (
	"github.com/theo1893/excelstreamer/data_source"
)

type SheetSelector struct {
	locator *SheetLocator
}

func NewSheetSelector() *SheetSelector {
	ins := &SheetSelector{
		locator: NewSheetLocator(),
	}

	return ins
}

// SelectSheet
// @Description: Select the sheet to which the dataSource belongs
func (s *SheetSelector) SelectSheet(dataSource data_source.DataSource, activeSheetList []SheetDefine) int {
	for index, sheet := range activeSheetList {
		if s.locator.Locate(dataSource, sheet.Anchor) {
			return index
		}
	}
	return -1
}
