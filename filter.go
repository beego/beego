package beego

import (
	"regexp"
)

type FilterRouter struct {
	pattern    string
	regex      *regexp.Regexp
	filterFunc FilterFunc
	hasregex   bool
}

func (mr *FilterRouter) ValidRouter(router string) bool {
	if mr.pattern == "" {
		return true
	}
	if mr.pattern == "*" {
		return true
	}
	if router == mr.pattern {
		return true
	}
	if mr.hasregex {
		if mr.regex.MatchString(router) {
			return true
		}
		matches := mr.regex.FindStringSubmatch(router)
		if len(matches[0]) == len(router) {
			return true
		}
	}
	return false
}
