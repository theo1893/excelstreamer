package excelstreamer

import (
	"github.com/theo1893/excelstreamer/consts"
	"github.com/theo1893/excelstreamer/data_source"
	"strconv"
)

func TransToStr(v interface{}) string {
	switch i := v.(type) {
	case string:
		return i
	case *string:
		return *i
	case int:
		return strconv.Itoa(i)
	case *int:
		return strconv.Itoa(*i)
	case int8:
		return strconv.Itoa(int(i))
	case *int8:
		return strconv.Itoa(int(*i))
	case int16:
		return strconv.Itoa(int(i))
	case *int16:
		return strconv.Itoa(int(*i))
	case int32:
		return strconv.Itoa(int(i))
	case *int32:
		return strconv.Itoa(int(*i))
	case int64:
		return strconv.Itoa(int(i))
	case *int64:
		return strconv.Itoa(int(*i))
	case float32:
		return strconv.FormatFloat(float64(i), 'f', 2, 64)
	case *float32:
		return strconv.FormatFloat(float64(*i), 'f', 2, 64)
	case float64:
		return strconv.FormatFloat(i, 'f', 2, 64)
	case *float64:
		return strconv.FormatFloat(*i, 'f', 2, 64)
	case data_source.Stringer:
		return v.(data_source.Stringer).String()
	default:

	}

	return consts.StringMinus
}
