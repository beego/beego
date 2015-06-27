// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package beego

import (
	"path"
	"regexp"
	"strings"

	"github.com/astaxie/beego/utils"
)

type Tree struct {
	//search fix route first
	fixrouters map[string]*Tree

	//if set, failure to match fixrouters search then search wildcard
	wildcard *Tree

	//if set, failure to match wildcard search
	leaves []*leafInfo
}

func NewTree() *Tree {
	return &Tree{
		fixrouters: make(map[string]*Tree),
	}
}

// add Tree to the exist Tree
// prefix should has no params
func (t *Tree) AddTree(prefix string, tree *Tree) {
	t.addtree(splitPath(prefix), tree, nil, "")
}

func (t *Tree) addtree(segments []string, tree *Tree, wildcards []string, reg string) {
	if len(segments) == 0 {
		panic("prefix should has path")
	}
	seg := segments[0]
	iswild, params, regexpStr := splitSegment(seg)
	if len(segments) == 1 {
		if iswild {
			if regexpStr != "" {
				if reg == "" {
					rr := ""
					for _, w := range wildcards {
						if w == "." || w == ":" {
							continue
						}
						if w == ":splat" {
							rr = rr + "(.+)/"
						} else {
							rr = rr + "([^/]+)/"
						}
					}
					regexpStr = rr + regexpStr
				} else {
					regexpStr = "/" + regexpStr
				}
			} else if reg != "" {
				if seg == "*.*" {
					regexpStr = "([^.]+).(.+)"
				} else {
					for _, w := range params {
						if w == "." || w == ":" {
							continue
						}
						regexpStr = "([^/]+)/" + regexpStr
					}
				}
			}
			reg = strings.Trim(reg+"/"+regexpStr, "/")
			filterTreeWithPrefix(tree, append(wildcards, params...), reg)
			t.wildcard = tree
		} else {
			reg = strings.Trim(reg+"/"+regexpStr, "/")
			filterTreeWithPrefix(tree, append(wildcards, params...), reg)
			t.fixrouters[seg] = tree
		}
		return
	}
	if iswild {
		if t.wildcard == nil {
			t.wildcard = NewTree()
		}
		if regexpStr != "" {
			if reg == "" {
				rr := ""
				for _, w := range wildcards {
					if w == "." || w == ":" {
						continue
					}
					if w == ":splat" {
						rr = rr + "(.+)/"
					} else {
						rr = rr + "([^/]+)/"
					}
				}
				regexpStr = rr + regexpStr
			} else {
				regexpStr = "/" + regexpStr
			}
		} else if reg != "" {
			if seg == "*.*" {
				regexpStr = "([^.]+).(.+)"
			} else {
				for _, w := range params {
					if w == "." || w == ":" {
						continue
					}
					regexpStr = "([^/]+)/" + regexpStr
				}
			}
		}
		reg = strings.TrimRight(strings.TrimRight(reg, "/")+"/"+regexpStr, "/")
		t.wildcard.addtree(segments[1:], tree, append(wildcards, params...), reg)
	} else {
		subTree := NewTree()
		t.fixrouters[seg] = subTree
		subTree.addtree(segments[1:], tree, append(wildcards, params...), reg)
	}
}

func filterTreeWithPrefix(t *Tree, wildcards []string, reg string) {
	for _, v := range t.fixrouters {
		filterTreeWithPrefix(v, wildcards, reg)
	}
	if t.wildcard != nil {
		filterTreeWithPrefix(t.wildcard, wildcards, reg)
	}
	for _, l := range t.leaves {
		if reg != "" {
			if l.regexps != nil {
				l.wildcards = append(wildcards, l.wildcards...)
				l.regexps = regexp.MustCompile("^" + reg + "/" + strings.Trim(l.regexps.String(), "^$") + "$")
			} else {
				for _, v := range l.wildcards {
					if v == ":" || v == "." {
						continue
					}
					if v == ":splat" {
						reg = reg + "/(.+)"
					} else {
						reg = reg + "/([^/]+)"
					}
				}
				l.regexps = regexp.MustCompile("^" + reg + "$")
				l.wildcards = append(wildcards, l.wildcards...)
			}
			filterCards := []string{}
			for _, v := range l.wildcards {
				if v == ":" || v == "." {
					continue
				}
				filterCards = append(filterCards, v)
			}
			l.wildcards = filterCards
		} else {
			l.wildcards = append(wildcards, l.wildcards...)
			if l.regexps != nil {
				for _, w := range wildcards {
					if w == "." || w == ":" {
						continue
					}
					if w == ":splat" {
						reg = "(.+)/" + reg
					} else {
						reg = "([^/]+)/" + reg
					}
				}
				l.regexps = regexp.MustCompile("^" + reg + strings.Trim(l.regexps.String(), "^$") + "$")
			}
		}
	}
}

// call addseg function
func (t *Tree) AddRouter(pattern string, runObject interface{}) {
	t.addseg(splitPath(pattern), runObject, nil, "")
}

// "/"
// "admin" ->
func (t *Tree) addseg(segments []string, route interface{}, wildcards []string, reg string) {
	if len(segments) == 0 {
		if reg != "" {
			filterCards := []string{}
			for _, v := range wildcards {
				if v == ":" || v == "." {
					continue
				}
				filterCards = append(filterCards, v)
			}
			t.leaves = append(t.leaves, &leafInfo{runObject: route, wildcards: filterCards, regexps: regexp.MustCompile("^" + reg + "$")})
		} else {
			t.leaves = append(t.leaves, &leafInfo{runObject: route, wildcards: wildcards})
		}
	} else {
		seg := segments[0]
		iswild, params, regexpStr := splitSegment(seg)
		//for the router  /login/*/access match /login/2009/11/access
		if !iswild && utils.InSlice(":splat", wildcards) {
			iswild = true
			regexpStr = seg
		}
		if seg == "*" && len(wildcards) > 0 && reg == "" {
			iswild = true
			regexpStr = "(.+)"
		}
		if iswild {
			if t.wildcard == nil {
				t.wildcard = NewTree()
			}
			if regexpStr != "" {
				if reg == "" {
					rr := ""
					for _, w := range wildcards {
						if w == "." || w == ":" {
							continue
						}
						if w == ":splat" {
							rr = rr + "(.+)/"
						} else {
							rr = rr + "([^/]+)/"
						}
					}
					regexpStr = rr + regexpStr
				} else {
					regexpStr = "/" + regexpStr
				}
			} else if reg != "" {
				if seg == "*.*" {
					regexpStr = "/([^.]+).(.+)"
				} else {
					for _, w := range params {

						if w == "." || w == ":" {
							continue
						}
						regexpStr = "/([^/]+)" + regexpStr
					}
				}
			}
			t.wildcard.addseg(segments[1:], route, append(wildcards, params...), reg+regexpStr)
		} else {
			subTree, ok := t.fixrouters[seg]
			if !ok {
				subTree = NewTree()
				t.fixrouters[seg] = subTree
			}
			subTree.addseg(segments[1:], route, wildcards, reg)
		}
	}
}

// match router to runObject & params
func (t *Tree) Match(pattern string) (runObject interface{}, params map[string]string) {
	if len(pattern) == 0 || pattern[0] != '/' {
		return nil, nil
	}

	return t.match(splitPath(pattern), nil)
}

func (t *Tree) match(segments []string, wildcardValues []string) (runObject interface{}, params map[string]string) {
	// Handle leaf nodes:
	if len(segments) == 0 {
		for _, l := range t.leaves {
			if ok, pa := l.match(wildcardValues); ok {
				return l.runObject, pa
			}
		}
		if t.wildcard != nil {
			for _, l := range t.wildcard.leaves {
				if ok, pa := l.match(wildcardValues); ok {
					return l.runObject, pa
				}
			}

		}
		return nil, nil
	}

	seg, segs := segments[0], segments[1:]

	subTree, ok := t.fixrouters[seg]
	if ok {
		runObject, params = subTree.match(segs, wildcardValues)
	} else if len(segs) == 0 { //.json .xml
		if subindex := strings.LastIndex(seg, "."); subindex != -1 {
			subTree, ok = t.fixrouters[seg[:subindex]]
			if ok {
				runObject, params = subTree.match(segs, wildcardValues)
				if runObject != nil {
					if params == nil {
						params = make(map[string]string)
					}
					params[":ext"] = seg[subindex+1:]
					return runObject, params
				}
			}
		}
	}
	if runObject == nil && t.wildcard != nil {
		runObject, params = t.wildcard.match(segs, append(wildcardValues, seg))
	}
	if runObject == nil {
		for _, l := range t.leaves {
			if ok, pa := l.match(append(wildcardValues, segments...)); ok {
				return l.runObject, pa
			}
		}
	}
	return runObject, params
}

type leafInfo struct {
	// names of wildcards that lead to this leaf. eg, ["id" "name"] for the wildcard ":id" and ":name"
	wildcards []string

	// if the leaf is regexp
	regexps *regexp.Regexp

	runObject interface{}
}

func (leaf *leafInfo) match(wildcardValues []string) (ok bool, params map[string]string) {
	if leaf.regexps == nil {
		// has error
		if len(wildcardValues) == 0 && len(leaf.wildcards) > 0 {
			if utils.InSlice(":", leaf.wildcards) {
				params = make(map[string]string)
				j := 0
				for _, v := range leaf.wildcards {
					if v == ":" {
						continue
					}
					params[v] = ""
					j += 1
				}
				return true, params
			}
			return false, nil
		} else if len(wildcardValues) == 0 { // static path
			return true, nil
		}
		// match *
		if len(leaf.wildcards) == 1 && leaf.wildcards[0] == ":splat" {
			params = make(map[string]string)
			params[":splat"] = path.Join(wildcardValues...)
			return true, params
		}
		// match *.*
		if len(leaf.wildcards) == 3 && leaf.wildcards[0] == "." {
			params = make(map[string]string)
			lastone := wildcardValues[len(wildcardValues)-1]
			strs := strings.SplitN(lastone, ".", 2)
			if len(strs) == 2 {
				params[":ext"] = strs[1]
			} else {
				params[":ext"] = ""
			}
			params[":path"] = path.Join(wildcardValues[:len(wildcardValues)-1]...) + "/" + strs[0]
			return true, params
		}
		// match :id
		params = make(map[string]string)
		j := 0
		for _, v := range leaf.wildcards {
			if v == ":" {
				continue
			}
			if v == "." {
				lastone := wildcardValues[len(wildcardValues)-1]
				strs := strings.SplitN(lastone, ".", 2)
				if len(strs) == 2 {
					params[":ext"] = strs[1]
				} else {
					params[":ext"] = ""
				}
				if len(wildcardValues[j:]) == 1 {
					params[":path"] = strs[0]
				} else {
					params[":path"] = path.Join(wildcardValues[j:]...) + "/" + strs[0]
				}
				return true, params
			}
			if len(wildcardValues) <= j {
				return false, nil
			}
			params[v] = wildcardValues[j]
			j += 1
		}
		if len(params) != len(wildcardValues) {
			return false, nil
		}
		return true, params
	}

	if !leaf.regexps.MatchString(path.Join(wildcardValues...)) {
		return false, nil
	}
	params = make(map[string]string)
	matches := leaf.regexps.FindStringSubmatch(path.Join(wildcardValues...))
	for i, match := range matches[1:] {
		params[leaf.wildcards[i]] = match
	}
	return true, params
}

// "/" -> []
// "/admin" -> ["admin"]
// "/admin/" -> ["admin"]
// "/admin/users" -> ["admin", "users"]
func splitPath(key string) []string {
	if key == "" {
		return []string{}
	}
	elements := strings.Split(key, "/")
	if elements[0] == "" {
		elements = elements[1:]
	}
	if elements[len(elements)-1] == "" {
		elements = elements[:len(elements)-1]
	}
	return elements
}

// "admin" -> false, nil, ""
// ":id" -> true, [:id], ""
// "?:id" -> true, [: :id], ""        : meaning can empty
// ":id:int" -> true, [:id], ([0-9]+)
// ":name:string" -> true, [:name], ([\w]+)
// ":id([0-9]+)" -> true, [:id], ([0-9]+)
// ":id([0-9]+)_:name" -> true, [:id :name], ([0-9]+)_(.+)
// "cms_:id_:page.html" -> true, [:id :page], cms_(.+)_(.+).html
// "*" -> true, [:splat], ""
// "*.*" -> true,[. :path :ext], ""      . meaning separator
func splitSegment(key string) (bool, []string, string) {
	if strings.HasPrefix(key, "*") {
		if key == "*.*" {
			return true, []string{".", ":path", ":ext"}, ""
		} else {
			return true, []string{":splat"}, ""
		}
	}
	if strings.ContainsAny(key, ":") {
		var paramsNum int
		var out []rune
		var start bool
		var startexp bool
		var param []rune
		var expt []rune
		var skipnum int
		params := []string{}
		reg := regexp.MustCompile(`[a-zA-Z0-9_]+`)
		for i, v := range key {
			if skipnum > 0 {
				skipnum -= 1
				continue
			}
			if start {
				//:id:int and :name:string
				if v == ':' {
					if len(key) >= i+4 {
						if key[i+1:i+4] == "int" {
							out = append(out, []rune("([0-9]+)")...)
							params = append(params, ":"+string(param))
							start = false
							startexp = false
							skipnum = 3
							param = make([]rune, 0)
							paramsNum += 1
							continue
						}
					}
					if len(key) >= i+7 {
						if key[i+1:i+7] == "string" {
							out = append(out, []rune(`([\w]+)`)...)
							params = append(params, ":"+string(param))
							paramsNum += 1
							start = false
							startexp = false
							skipnum = 6
							param = make([]rune, 0)
							continue
						}
					}
				}
				// params only support a-zA-Z0-9
				if reg.MatchString(string(v)) {
					param = append(param, v)
					continue
				}
				if v != '(' {
					out = append(out, []rune(`(.+)`)...)
					params = append(params, ":"+string(param))
					param = make([]rune, 0)
					paramsNum += 1
					start = false
					startexp = false
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
				start = true
			} else if v == '(' {
				startexp = true
				start = false
				params = append(params, ":"+string(param))
				paramsNum += 1
				expt = make([]rune, 0)
				expt = append(expt, '(')
			} else if v == ')' {
				startexp = false
				expt = append(expt, ')')
				out = append(out, expt...)
				param = make([]rune, 0)
			} else if v == '?' {
				params = append(params, ":")
			} else {
				out = append(out, v)
			}
		}
		if len(param) > 0 {
			if paramsNum > 0 {
				out = append(out, []rune(`(.+)`)...)
			}
			params = append(params, ":"+string(param))
		}
		return true, params, string(out)
	} else {
		return false, nil, ""
	}
}
