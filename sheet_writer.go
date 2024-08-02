package excelstreamer

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

type SheetWriter struct {
	sw            *excelize.StreamWriter
	rawFile       *excelize.File
	sheetName     string // basic sheet name
	sheetIndex    int    // current sheet index
	overflowLimit int
	nextRow       int
}

type Field struct {
	Header      string
	HeaderStyle *excelize.Style
	DataStyle   *excelize.Style
	Data        interface{}
	ColWidth    float64
	RowHeight   float64
}

func newSheetWriter(file *excelize.File, sheetName string, overflowLimit int) *SheetWriter {
	_, err := file.NewSheet(sheetName)
	if err != nil {
		return nil
	}

	sw, err := file.NewStreamWriter(sheetName)
	if err != nil {
		return nil
	}

	return &SheetWriter{
		sw:            sw,
		rawFile:       file,
		sheetName:     sheetName,
		overflowLimit: overflowLimit,
		nextRow:       1,
		sheetIndex:    0,
	}
}

// WriteHeaderToSheet
// @Description: Write only the headers to sheet
func (s *SheetWriter) WriteHeaderToSheet(fieldList []Field, styleId int) error {
	if len(fieldList) == 0 {
	}

	if s.Dirty() {
		return nil
	}

	headerList := make([]interface{}, len(fieldList))
	colWidthList := make([]float64, len(fieldList))

	for index := range fieldList {
		headerList[index] = excelize.Cell{
			StyleID: styleId,
			Value:   fieldList[index].Header,
		}

		if fieldList[index].HeaderStyle != nil {
			customStyleId, _ := s.rawFile.NewStyle(fieldList[index].HeaderStyle)
			headerList[index] = excelize.Cell{
				StyleID: customStyleId,
				Value:   fieldList[index].Header,
			}
		} else if index > 0 && fieldList[index-1].HeaderStyle != nil {
			fieldList[index].HeaderStyle = fieldList[index-1].HeaderStyle
			customStyleId, _ := s.rawFile.NewStyle(fieldList[index].HeaderStyle)
			headerList[index] = excelize.Cell{
				StyleID: customStyleId,
				Value:   fieldList[index].Header,
			}
		}

		colWidthList[index] = fieldList[index].ColWidth
	}

	for colIndex, colWidth := range colWidthList {
		setErr := s.sw.SetColWidth(colIndex+1, colIndex+1, colWidth)
		if setErr != nil {

		}
	}

	writeErr := s.sw.SetRow(s.CurAxis(), headerList, excelize.RowOpts{
		Height: fieldList[0].RowHeight,
		Hidden: false,
	})
	if writeErr != nil {

	}
	s.Populate()

	return nil
}

// WriteDataToSheet
// @Description: Write data to sheet
func (s *SheetWriter) WriteDataToSheet(fieldList []Field, styleId int) error {
	if len(fieldList) == 0 {
	}

	// If the sheet has overflowed, a new sheet will be created to accommodate following data
	if s.Overflow() {
		flushErr := s.sw.Flush()
		if flushErr != nil {
			return flushErr
		}

		_, _ = s.rawFile.NewSheet(s.NextSheetName())
		newSw, createErr := s.rawFile.NewStreamWriter(s.NextSheetName())
		s.Reset()

		if createErr != nil {

		}

		s.sw = newSw
	}

	headerList := make([]interface{}, len(fieldList))
	dataList := make([]interface{}, len(fieldList))
	colWidthList := make([]float64, len(fieldList))

	for index := range fieldList {
		if !s.Dirty() {
			headerList[index] = excelize.Cell{
				StyleID: styleId,
				Value:   fieldList[index].Header,
			}

			if fieldList[index].HeaderStyle != nil {
				customStyleId, _ := s.rawFile.NewStyle(fieldList[index].HeaderStyle)
				headerList[index] = excelize.Cell{
					StyleID: customStyleId,
					Value:   fieldList[index].Header,
				}
			} else if index > 0 && fieldList[index-1].HeaderStyle != nil {
				fieldList[index].HeaderStyle = fieldList[index-1].HeaderStyle
				customStyleId, _ := s.rawFile.NewStyle(fieldList[index].HeaderStyle)
				headerList[index] = excelize.Cell{
					StyleID: customStyleId,
					Value:   fieldList[index].Header,
				}
			}
			colWidthList[index] = fieldList[index].ColWidth
		}

		dataList[index] = excelize.Cell{
			StyleID: styleId,
			Value:   fieldList[index].Data,
		}
		if fieldList[index].DataStyle != nil {
			customStyleId, _ := s.rawFile.NewStyle(fieldList[index].DataStyle)
			dataList[index] = excelize.Cell{
				StyleID: customStyleId,
				Value:   fieldList[index].Data,
			}
		} else if index > 0 && fieldList[index-1].DataStyle != nil {
			fieldList[index].DataStyle = fieldList[index-1].DataStyle
			customStyleId, _ := s.rawFile.NewStyle(fieldList[index].DataStyle)
			dataList[index] = excelize.Cell{
				StyleID: customStyleId,
				Value:   fieldList[index].Data,
			}
		}
	}

	if !s.Dirty() {
		for colIndex, colWidth := range colWidthList {
			setErr := s.sw.SetColWidth(colIndex+1, colIndex+1, colWidth)
			if setErr != nil {
			}
		}

		writeErr := s.sw.SetRow(s.CurAxis(), headerList, excelize.RowOpts{
			Height: fieldList[0].RowHeight,
			Hidden: false,
		})
		if writeErr != nil {

		}
		s.Populate()
	}

	writeErr := s.sw.SetRow(s.CurAxis(), dataList, excelize.RowOpts{
		Height: fieldList[0].RowHeight,
		Hidden: false,
	})
	if writeErr != nil {

	}
	s.Populate()

	return nil
}

func (s *SheetWriter) Flush() error {
	if flushErr := s.sw.Flush(); flushErr != nil {
		return flushErr
	}

	return nil
}

func (s *SheetWriter) CurAxis() string {
	return fmt.Sprintf("A%d", s.nextRow)
}

func (s *SheetWriter) Reset() {
	s.sheetIndex++
	s.nextRow = 1
}

func (s *SheetWriter) Populate() {
	s.nextRow++
}

func (s *SheetWriter) Dirty() bool {
	return s.nextRow > 1
}

func (s *SheetWriter) Overflow() bool {
	return s.nextRow > s.overflowLimit
}

func (s *SheetWriter) NextSheetName() string {
	return fmt.Sprintf("%s_%d", s.sheetName, s.sheetIndex+1)
}
