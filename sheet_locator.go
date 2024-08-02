package excelstreamer

import (
	"github.com/theo1893/excelstreamer/data_source"
	"regexp"
	"strings"
)

type SheetLocator struct {
	matcher *regexp.Regexp
}

func NewSheetLocator() *SheetLocator {
	return &SheetLocator{
		// example: name=Alex&&age=24||name=Ego
		matcher: regexp.MustCompile(`(([\w.]+)=([\w]+))(&&(([\w.]+)=([\w]+)))*`),
	}
}

func (s *SheetLocator) Locate(dataSource data_source.DataSource, anchor string) bool {
	result := s.matcher.FindAllStringSubmatch(anchor, -1)

	if len(result) == 0 {
		return false
	}

	for _, item := range result {
		conditions := strings.Split(item[0], "&&")
		pass := true
		for _, condition := range conditions {
			kv := strings.Split(condition, "=")
			keys := strings.Split(kv[0], ".")
			if TransToStr(dataSource.GetInterface(keys...)) != kv[1] {
				pass = false
				break
			}
		}

		if pass {
			return true
		}
	}
	return false
}
