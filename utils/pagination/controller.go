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

// Usage
//
// In your beego.Controller:
//
//  package controllers
//
//  import "github.com/astaxie/beego/utils/pagination"
//
//  type PostsController struct {
//    beego.Controller
//  }
//
//  func (this *PostsController) ListAllPosts() {
//      // sets this.Data["paginator"] with the current offset (from the url query param)
//      postsPerPage := 20
//      paginator := pagination.SetPaginator(this, postsPerPage, CountPosts())
//
//      // fetch the next 20 posts
//      this.Data["posts"] = ListPostsByOffsetAndLimit(paginator.Offset(), postsPerPage)
//  }
//
//
// In your view templates:
//
//  {{if .paginator.HasPages}}
//  <ul class="pagination pagination">
//      {{if .paginator.HasPrev}}
//          <li><a href="{{.paginator.PageLinkFirst}}">{{ i18n .Lang "paginator.first_page"}}</a></li>
//          <li><a href="{{.paginator.PageLinkPrev}}">&laquo;</a></li>
//      {{else}}
//          <li class="disabled"><a>{{ i18n .Lang "paginator.first_page"}}</a></li>
//          <li class="disabled"><a>&laquo;</a></li>
//      {{end}}
//      {{range $index, $page := .paginator.Pages}}
//          <li{{if $.paginator.IsActive .}} class="active"{{end}}>
//              <a href="{{$.paginator.PageLink $page}}">{{$page}}</a>
//          </li>
//      {{end}}
//      {{if .paginator.HasNext}}
//          <li><a href="{{.paginator.PageLinkNext}}">&raquo;</a></li>
//          <li><a href="{{.paginator.PageLinkLast}}">{{ i18n .Lang "paginator.last_page"}}</a></li>
//      {{else}}
//          <li class="disabled"><a>&raquo;</a></li>
//          <li class="disabled"><a>{{ i18n .Lang "paginator.last_page"}}</a></li>
//      {{end}}
//  </ul>
//  {{end}}
//
// See also http://beego.me/docs/mvc/view/page.md
package pagination

import (
	"github.com/astaxie/beego/context"
)

type PaginationController interface {
	GetCtx() *context.Context
	GetData() map[interface{}]interface{}
}

// Instantiates a Paginator and assigns it to controller.Data["paginator"].
func SetPaginator(controller PaginationController, per int, nums int64) (paginator *Paginator) {
	request := controller.GetCtx().Request
	paginator = NewPaginator(request, per, nums)
	data := controller.GetData()
	data["paginator"] = paginator
	return
}
