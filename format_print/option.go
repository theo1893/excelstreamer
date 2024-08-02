package format_print

const (
	// format option
	FormatTIMESTAMP FormatOption = "TIMESTAMP"
)

var FormatFuncCollection = map[FormatOption]FormatFunc{
	FormatTIMESTAMP: FormatTIMESTAMPFunc,
}

type FormatFunc func(string, ...string) string
type FormatOption string

func RegisterFormat(option FormatOption, fc FormatFunc) {
	FormatFuncCollection[option] = fc
}
