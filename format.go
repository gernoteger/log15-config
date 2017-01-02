package config

import (
	"fmt"

	"github.com/inconshreveable/log15"
)

type Fmt string

// NewFormatFunc creates a new format
type NewFormatFunc func() log15.Format

var formats = map[string]func() log15.Format{
	"terminal": log15.TerminalFormat,
	"json":     log15.JsonFormat,
	"logfmt":   log15.LogfmtFormat,
}

// AddFormat adds a Format to the list. You can even replace the old ones!!
func AddFormat(key string, newFunc NewFormatFunc) {
	formats[key] = newFunc
}

//
//
//// FmtFromString returns the appropriate format string or errors out if unknown
//func FmtFromString(fmtString string) (Fmt, error) {
//
//	switch strings.ToLower(fmtString) {
//	case "terminal", "term", "console":
//		return FmtTerminal, nil
//	case "json":
//		return FmtJson, nil
//	case "logfmt":
//		return FmtLogfmt, nil
//	default:
//		return FmtTerminal, fmt.Errorf("Unknown format: %v", fmtString)
//	}
//}
//
//// UnmarshalString to implement StringUnmarshaller
//func (f Fmt) UnmarshalString(from string) (interface{}, error) {
//	f1, err := FmtFromString(from)
//	if err != nil {
//		return nil, err
//	}
//
//	return f1, nil
//}

func (f Fmt) NewFormat() log15.Format {

	if f == "" {
		f = "logfmt" // default
	}
	newFmt, ok := formats[string(f)]
	if !ok {
		err := fmt.Errorf("unknown format: '%v'", f)
		panic(err) //TODO: return errors??
	}

	return newFmt()
}
