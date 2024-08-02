package excelstreamer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/theo1893/excelstreamer/data_source"
	"github.com/theo1893/excelstreamer/streamer"
	"github.com/xuri/excelize/v2"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
)

func NewExcelStreamWriter(template string, confFile string, dstFile string) (streamer.Streamer, error) {
	var file *excelize.File
	if template == "" {
		file = excelize.NewFile()
	} else {
		f, err := excelize.OpenFile(template)
		if err != nil {
			return nil, fmt.Errorf("open excel failed, err=[%+v]", err.Error())
		}
		file = f
	}

	rawConf, err := os.ReadFile(confFile)
	if err != nil {
		return nil, fmt.Errorf("read config file failed, err=[%+v]", err.Error())
	}

	conf := ExcelDefine{}
	err = yaml.Unmarshal(rawConf, &conf)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config file failed, err=[%+v]", err.Error())
	}

	ctx := context.Background()
	s := streamer.NewStreamer(ctx)
	sheetSelector := NewSheetSelector()
	sheetLoader := NewSheetLoader()
	excelStreamWriter := NewExcelWriter(file, 1000000)

	activeSheetList := sheetLoader.FetchActiveSheet(conf)

	// transfer string to map
	s.AddKnot(func(i interface{}) (interface{}, error) {
		if _, ok := i.(string); !ok {
			return nil, fmt.Errorf("input is not string, %+v", i)
		}

		v := i.(string)
		m := make(map[string]interface{})

		d := json.NewDecoder(strings.NewReader(v))
		d.UseNumber()
		if e := d.Decode(&m); e != nil {
			return nil, fmt.Errorf("unmarshal json failed, err=[%+v]", e)
		}

		return data_source.MapDataSource(m), nil
	})

	// write data to sheet
	s.AddKnot(func(i interface{}) (interface{}, error) {
		if _, ok := i.(data_source.DataSource); !ok {
			return nil, fmt.Errorf("input is not DataSource, %+v", i)
		}

		dataSource := i.(data_source.DataSource)
		sheetIdx := sheetSelector.SelectSheet(dataSource, activeSheetList)
		if sheetIdx < 0 {
			return nil, nil
		}
		fieldList := sheetLoader.LoadData(dataSource, conf.Sheets[sheetIdx])
		e := excelStreamWriter.WriteDataToSheet(conf.Sheets[sheetIdx].SheetName, fieldList, 0)
		return nil, e
	})

	// add flush as release func
	s.AddReleaseFunc(func() {
		_ = excelStreamWriter.Flush()
		_ = file.SaveAs(dstFile)
		_ = file.Close()
		log.Println("Release func is called.")
	})

	return s, nil
}
