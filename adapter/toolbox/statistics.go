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

package toolbox

import (
	"time"

	"github.com/astaxie/beego/server/web"
)

// Statistics struct
type Statistics web.Statistics

// URLMap contains several statistics struct to log different data
type URLMap web.URLMap

// AddStatistics add statistics task.
// it needs request method, request url, request controller and statistics time duration
func (m *URLMap) AddStatistics(requestMethod, requestURL, requestController string, requesttime time.Duration) {
	(*web.URLMap)(m).AddStatistics(requestMethod, requestURL, requestController, requesttime)
}

// GetMap put url statistics result in io.Writer
func (m *URLMap) GetMap() map[string]interface{} {
	return (*web.URLMap)(m).GetMap()
}

// GetMapData return all mapdata
func (m *URLMap) GetMapData() []map[string]interface{} {
	return (*web.URLMap)(m).GetMapData()
}

// StatisticsMap hosld global statistics data map
var StatisticsMap *URLMap

func init() {
	StatisticsMap = (*URLMap)(web.StatisticsMap)
}
