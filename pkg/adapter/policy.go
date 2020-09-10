// Copyright 2016 beego authors. All Rights Reserved.
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

package adapter

import (
	"github.com/astaxie/beego/pkg/adapter/context"
	"github.com/astaxie/beego/pkg/server/web"
	beecontext "github.com/astaxie/beego/pkg/server/web/context"
)

// PolicyFunc defines a policy function which is invoked before the controller handler is executed.
type PolicyFunc func(*context.Context)

// FindPolicy Find Router info for URL
func (p *ControllerRegister) FindPolicy(cont *context.Context) []PolicyFunc {
	pf := (*web.ControllerRegister)(p).FindPolicy((*beecontext.Context)(cont))
	npf := newToOldPolicyFunc(pf)
	return npf
}

func newToOldPolicyFunc(pf []web.PolicyFunc) []PolicyFunc {
	npf := make([]PolicyFunc, 0, len(pf))
	for _, f := range pf {
		npf = append(npf, func(c *context.Context) {
			f((*beecontext.Context)(c))
		})
	}
	return npf
}

func oldToNewPolicyFunc(pf []PolicyFunc) []web.PolicyFunc {
	npf := make([]web.PolicyFunc, 0, len(pf))
	for _, f := range pf {
		npf = append(npf, func(c *beecontext.Context) {
			f((*context.Context)(c))
		})
	}
	return npf
}

// Policy Register new policy in beego
func Policy(pattern, method string, policy ...PolicyFunc) {
	pf := oldToNewPolicyFunc(policy)
	web.Policy(pattern, method, pf...)
}
