package parser

import (
	"strings"
	"time"

	"github.com/pingcap/tidb/infoschema/inspection/log/item"
)

const (
	TimeStampLayout       = "2006/01/02 15:04:05.000 -07:00"
	FormerTimeStampLayout = "2006/01/02 15:04:05.000"
)

var LevelTypeMap = map[string]item.LevelType{
	"CRITICAL": item.LevelFATAL,
	"FATAL":    item.LevelFATAL,
	"ERROR":    item.LevelERROR,
	"ERRO":     item.LevelERROR,
	"WARNING":  item.LevelWARN,
	"WARN":     item.LevelWARN,
	"INFO":     item.LevelINFO,
	"DEBUG":    item.LevelDEBUG,
}

func ParseLogLevel(b []byte) item.LevelType {
	s := strings.ToUpper(string(b))
	if level, ok := LevelTypeMap[s]; ok {
		return level
	} else {
		return item.LevelInvalid
	}
}

// TiDB / TiKV / PD unified log format
// [2019/03/04 17:04:24.614 +08:00] ...
func parseTimeStamp(b []byte) (*time.Time, error) {
	t, err := time.Parse(TimeStampLayout, string(b))
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// TiDB / TiKV / PD log format used in former version
// 2019/07/18 11:04:29.314 ...
func parseFormerTimeStamp(b []byte) (*time.Time, error) {
	local, err := time.LoadLocation("Asia/Chongqing")
	if err != nil {
		return nil, err
	}
	t, err := time.ParseInLocation(FormerTimeStampLayout, string(b), local)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
