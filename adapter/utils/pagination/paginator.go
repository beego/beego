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

package pagination

import (
	"net/http"

	"github.com/beego/beego/core/utils/pagination"
)

// Paginator within the state of a http request.
type Paginator pagination.Paginator

// PageNums Returns the total number of pages.
func (p *Paginator) PageNums() int {
	return (*pagination.Paginator)(p).PageNums()
}

// Nums Returns the total number of items (e.g. from doing SQL count).
func (p *Paginator) Nums() int64 {
	return (*pagination.Paginator)(p).Nums()
}

// SetNums Sets the total number of items.
func (p *Paginator) SetNums(nums interface{}) {
	(*pagination.Paginator)(p).SetNums(nums)
}

// Page Returns the current page.
func (p *Paginator) Page() int {
	return (*pagination.Paginator)(p).Page()
}

// Pages Returns a list of all pages.
//
// Usage (in a view template):
//
//  {{range $index, $page := .paginator.Pages}}
//    <li{{if $.paginator.IsActive .}} class="active"{{end}}>
//      <a href="{{$.paginator.PageLink $page}}">{{$page}}</a>
//    </li>
//  {{end}}
func (p *Paginator) Pages() []int {
	return (*pagination.Paginator)(p).Pages()
}

// PageLink Returns URL for a given page index.
func (p *Paginator) PageLink(page int) string {
	return (*pagination.Paginator)(p).PageLink(page)
}

// PageLinkPrev Returns URL to the previous page.
func (p *Paginator) PageLinkPrev() (link string) {
	return (*pagination.Paginator)(p).PageLinkPrev()
}

// PageLinkNext Returns URL to the next page.
func (p *Paginator) PageLinkNext() (link string) {
	return (*pagination.Paginator)(p).PageLinkNext()
}

// PageLinkFirst Returns URL to the first page.
func (p *Paginator) PageLinkFirst() (link string) {
	return (*pagination.Paginator)(p).PageLinkFirst()
}

// PageLinkLast Returns URL to the last page.
func (p *Paginator) PageLinkLast() (link string) {
	return (*pagination.Paginator)(p).PageLinkLast()
}

// HasPrev Returns true if the current page has a predecessor.
func (p *Paginator) HasPrev() bool {
	return (*pagination.Paginator)(p).HasPrev()
}

// HasNext Returns true if the current page has a successor.
func (p *Paginator) HasNext() bool {
	return (*pagination.Paginator)(p).HasNext()
}

// IsActive Returns true if the given page index points to the current page.
func (p *Paginator) IsActive(page int) bool {
	return (*pagination.Paginator)(p).IsActive(page)
}

// Offset Returns the current offset.
func (p *Paginator) Offset() int {
	return (*pagination.Paginator)(p).Offset()
}

// HasPages Returns true if there is more than one page.
func (p *Paginator) HasPages() bool {
	return (*pagination.Paginator)(p).HasPages()
}

// NewPaginator Instantiates a paginator struct for the current http request.
func NewPaginator(req *http.Request, per int, nums interface{}) *Paginator {
	return (*Paginator)(pagination.NewPaginator(req, per, nums))
}
