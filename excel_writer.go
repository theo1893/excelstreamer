package excelstreamer

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

type ExcelWriter struct {
	file          *excelize.File
	overflowLimit int // if the number of rows within one sheet is over this limit, a new sheet will be created
	sheetWriters  map[string]*SheetWriter
}

func NewExcelWriter(file *excelize.File, overflowLimit int) *ExcelWriter {
	return &ExcelWriter{
		file:          file,
		overflowLimit: overflowLimit,
		sheetWriters:  map[string]*SheetWriter{},
	}
}

func (s *ExcelWriter) WriteHeaderToSheet(sheetName string, fieldList []Field, styleId int) error {
	if _, ok := s.sheetWriters[sheetName]; !ok {
		s.sheetWriters[sheetName] = newSheetWriter(s.file, sheetName, s.overflowLimit)
	}

	sw := s.sheetWriters[sheetName]
	if sw == nil {
		return fmt.Errorf("init sheet writer failed")
	}

	return sw.WriteHeaderToSheet(fieldList, styleId)
}

func (s *ExcelWriter) WriteDataToSheet(sheetName string, fieldList []Field, styleId int) error {
	if _, ok := s.sheetWriters[sheetName]; !ok {
		s.sheetWriters[sheetName] = newSheetWriter(s.file, sheetName, s.overflowLimit)
	}

	sw := s.sheetWriters[sheetName]
	if sw == nil {
		return fmt.Errorf("init sheet writer failed")
	}

	return sw.WriteDataToSheet(fieldList, styleId)
}

func (s *ExcelWriter) Flush() error {
	for _, sw := range s.sheetWriters {
		if flushErr := sw.Flush(); flushErr != nil {
			return flushErr
		}
	}

	return nil
}
