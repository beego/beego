package goyaml2

import (
	"regexp"
	"strconv"
	//"time"
)

var (
	RE_INT, _   = regexp.Compile("^[0-9,]+$")
	RE_FLOAT, _ = regexp.Compile("^[0-9]+[.][0-9]+$")
	RE_DATE, _  = regexp.Compile("^[0-9]{4}-[0-9]{2}-[0-9]{2}$")
	RE_TIME, _  = regexp.Compile("^[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}$")
)

func string2Val(str string) interface{} {
	tmp := []byte(str)
	switch {
	case str == "false":
		return false
	case str == "true":
		return true
	case RE_INT.Match(tmp):
		// TODO check err
		_int, _ := strconv.ParseInt(str, 10, 64)
		return _int
	case RE_FLOAT.Match(tmp):
		_float, _ := strconv.ParseFloat(str, 64)
		return _float
		//TODO support time or Not?
		/*
			case RE_DATE.Match(tmp):
				_date, _ := time.Parse("2006-01-02", str)
				return _date
			case RE_TIME.Match(tmp):
				_time, _ := time.Parse("2006-01-02 03:04:05", str)
				return _time
		*/
	}
	return str
}
