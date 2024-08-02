package excelstreamer

import (
	"encoding/json"
	"github.com/theo1893/excelstreamer/consts"
	"github.com/theo1893/excelstreamer/data_source"
	"github.com/theo1893/excelstreamer/format_print"
	"github.com/xuri/excelize/v2"
	"strings"
)

type FieldLoader struct {
}

func NewFieldLoader() *FieldLoader {
	return &FieldLoader{}
}

// FieldLoad
// @Description: Extract data from dataSource to show in excel according to fieldDefine
func (f *FieldLoader) FieldLoad(dataSource data_source.DataSource, fieldDefine FieldDefine) Field {
	var format, formatArg string

	parsedValueSource := strings.Split(fieldDefine.Source, ".")
	stringData := ""

	// TODO Do we need cache?
	if fieldDefine.Format != "" {
		colonIdx := strings.Index(fieldDefine.Format, ":")
		if colonIdx != -1 {
			format = fieldDefine.Format[:colonIdx]
			formatArg = fieldDefine.Format[colonIdx+1:]
		} else {
			format = fieldDefine.Format
		}
	}

	if format_print.FormatFuncCollection[format_print.FormatOption(format)] != nil {
		stringData = format_print.FormatFuncCollection[format_print.FormatOption(format)](TransToStr(dataSource.GetInterface(parsedValueSource...)), strings.Split(formatArg, ";")...)
	} else {
		stringData = TransToStr(dataSource.GetInterface(parsedValueSource...))
	}

	var headerStyle, dataStyle *excelize.Style
	if fieldDefine.HeaderStyle != "" {
		headerStyle = &excelize.Style{}
		_ = json.Unmarshal([]byte(fieldDefine.HeaderStyle), headerStyle)
	}

	if fieldDefine.DataStyle != "" {
		dataStyle = &excelize.Style{}
		_ = json.Unmarshal([]byte(fieldDefine.DataStyle), dataStyle)
	}

	field := Field{
		ColWidth:    consts.DefaultSheetColWidth,
		RowHeight:   consts.DefaultSheetRowHeight,
		Data:        stringData,
		Header:      fieldDefine.Header,
		HeaderStyle: headerStyle,
		DataStyle:   dataStyle,
	}

	if fieldDefine.Width != 0 {
		field.ColWidth = float64(fieldDefine.Width)
	}

	return field
}
