package beego

import (
	"regexp"
	"strings"
)

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

func buildFilter(pattern string, filter FilterFunc) *FilterRouter {
	mr := new(FilterRouter)
	mr.params = make(map[int]string)
	mr.filterFunc = filter
	parts := strings.Split(pattern, "/")
	j := 0
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			expr := "(.+)"
			//a user may choose to override the default expression
			// similar to expressjs: ‘/user/:id([0-9]+)’
			if index := strings.Index(part, "("); index != -1 {
				expr = part[index:]
				part = part[:index]
				//match /user/:id:int ([0-9]+)
				//match /post/:username:string	([\w]+)
			} else if lindex := strings.LastIndex(part, ":"); lindex != 0 {
				switch part[lindex:] {
				case ":int":
					expr = "([0-9]+)"
					part = part[:lindex]
				case ":string":
					expr = `([\w]+)`
					part = part[:lindex]
				}
			}
			mr.params[j] = part
			parts[i] = expr
			j++
		}
		if strings.HasPrefix(part, "*") {
			expr := "(.+)"
			if part == "*.*" {
				mr.params[j] = ":path"
				parts[i] = "([^.]+).([^.]+)"
				j++
				mr.params[j] = ":ext"
				j++
			} else {
				mr.params[j] = ":splat"
				parts[i] = expr
				j++
			}
		}
		//url like someprefix:id(xxx).html
		if strings.Contains(part, ":") && strings.Contains(part, "(") && strings.Contains(part, ")") {
			var out []rune
			var start bool
			var startexp bool
			var param []rune
			var expt []rune
			for _, v := range part {
				if start {
					if v != '(' {
						param = append(param, v)
						continue
					}
				}
				if startexp {
					if v != ')' {
						expt = append(expt, v)
						continue
					}
				}
				if v == ':' {
					param = make([]rune, 0)
					param = append(param, ':')
					start = true
				} else if v == '(' {
					startexp = true
					start = false
					mr.params[j] = string(param)
					j++
					expt = make([]rune, 0)
					expt = append(expt, '(')
				} else if v == ')' {
					startexp = false
					expt = append(expt, ')')
					out = append(out, expt...)
				} else {
					out = append(out, v)
				}
			}
			parts[i] = string(out)
		}
	}

	if j != 0 {
		pattern = strings.Join(parts, "/")
		regex, regexErr := regexp.Compile(pattern)
		if regexErr != nil {
			//TODO add error handling here to avoid panic
			panic(regexErr)
		}
		mr.regex = regex
		mr.hasregex = true
	}
	mr.pattern = pattern
	return mr
}
