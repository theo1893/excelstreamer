package excelstreamer

import (
	"github.com/theo1893/excelstreamer/data_source"
)

type SheetLoader struct {
	fieldLoader *FieldLoader
}

func NewSheetLoader() *SheetLoader {
	return &SheetLoader{
		fieldLoader: NewFieldLoader(),
	}
}

// LoadData
// @Description: Construct excel cells according to the sheet define and data source
func (l *SheetLoader) LoadData(dataSource data_source.DataSource, sheet SheetDefine) []Field {
	sheetDataList := make([]Field, len(sheet.Fields))

	for idx, fieldDefine := range sheet.Fields {
		sheetDataList[idx] = l.fieldLoader.FieldLoad(dataSource, fieldDefine)
	}

	return sheetDataList
}

func (l *SheetLoader) LoadHeader(activeSheetList []SheetDefine) [][]Field {
	sheetHeaderList := make([][]Field, 0, len(activeSheetList))

	for _, sheet := range activeSheetList {
		fooDataSource := data_source.NewMapDataSource(make(map[string]interface{}))
		sheetDataList := l.LoadData(fooDataSource, sheet)
		sheetHeaderList = append(sheetHeaderList, sheetDataList)
	}

	return sheetHeaderList
}

func (l *SheetLoader) FetchActiveSheet(template ExcelDefine) []SheetDefine {
	activeSheetList := make([]SheetDefine, 0, len(template.Sheets))
	for _, sheet := range template.Sheets {
		if sheet.Hide {
			continue
		}
		activeSheetList = append(activeSheetList, sheet)
	}

	return activeSheetList
}
