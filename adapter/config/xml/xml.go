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

// Package xml for config provider.
//
// depend on github.com/beego/x2j.
//
// go install github.com/beego/x2j.
//
// Usage:
//
//	import(
//	  _ "github.com/beego/beego/v2/core/config/xml"
//	    "github.com/beego/beego/v2/core/config"
//	)
//
//	cnf, err := config.NewConfig("xml", "config.xml")
package xml

import (
	_ "github.com/beego/beego/v2/core/config/xml"
)
