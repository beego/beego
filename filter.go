// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie

package beego

// FilterRouter defines filter operation before controller handler execution.
// it can match patterned url and do filter function when action arrives.
type FilterRouter struct {
	filterFunc FilterFunc
	tree       *Tree
	pattern    string
}

// ValidRouter check current request is valid for this filter.
// if matched, returns parsed params in this request by defined filter router pattern.
func (f *FilterRouter) ValidRouter(router string) (bool, map[string]string) {
	isok, params := f.tree.Match(router)
	if isok == nil {
		return false, nil
	}
	if isok, ok := isok.(bool); ok {
		return isok, params
	} else {
		return false, nil
	}
}
