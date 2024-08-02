package format_print

import (
	"github.com/theo1893/excelstreamer/consts"
	"strconv"
	"time"
)

// FormatTIMESTAMPFunc
// @Description: Format timestamp as template
func FormatTIMESTAMPFunc(rawData string, args ...string) string {
	timestamp, err := strconv.ParseInt(rawData, 10, 64)
	if err != nil {
		return consts.StringMinus
	}

	if len(args) != 1 {
		return consts.StringMinus
	}

	return time.Unix(timestamp, 0).Format(args[0])
}
