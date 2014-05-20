// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/astaxie/beego for the canonical source repository
// @license     http://github.com/astaxie/beego/blob/master/LICENSE
// @authors     astaxie

package beego

import "regexp"

// FilterRouter defines filter operation before controller handler execution.
// it can match patterned url and do filter function when action arrives.
type FilterRouter struct {
	pattern     string
	regex       *regexp.Regexp
	filterFunc  FilterFunc
	hasregex    bool
	params      map[int]string
	parseParams map[string]string
}

// ValidRouter check current request is valid for this filter.
// if matched, returns parsed params in this request by defined filter router pattern.
func (mr *FilterRouter) ValidRouter(router string) (bool, map[string]string) {
	if mr.pattern == "" {
		return true, nil
	}
	if mr.pattern == "*" {
		return true, nil
	}
	if router == mr.pattern {
		return true, nil
	}
	//pattern /admin  router /admin/  match
	//pattern /admin/ router /admin don't match, because url will 301 in router
	if n := len(router); n > 1 && router[n-1] == '/' && router[:n-2] == mr.pattern {
		return true, nil
	}

	if mr.hasregex {
		if !mr.regex.MatchString(router) {
			return false, nil
		}
		matches := mr.regex.FindStringSubmatch(router)
		if len(matches) > 0 {
			if len(matches[0]) == len(router) {
				params := make(map[string]string)
				for i, match := range matches[1:] {
					params[mr.params[i]] = match
				}
				return true, params
			}
		}
	}
	return false, nil
}
