package utils

import (
	"golang.org/x/text/encoding/simplifiedchinese"
	"regexp"
)

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

// ConvertByte2String 解码
func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}

func GetValueStringByRegex(str, rule string) []string {
	reg, err := regexp.Compile(rule)
	if reg == nil || err != nil {
		return nil
	}
	//提取关键信息
	result := reg.FindStringSubmatch(str)
	if len(result) < 2 {
		return nil
	}
	return result
}
