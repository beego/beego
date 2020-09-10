// Copyright 2020
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bean

import (
	"context"
)

// TypeAdapter is an abstraction that define some behavior of target type
// usually, we don't use this to support basic type since golang has many restriction for basic types
// This is an important extension point
type TypeAdapter interface {
	DefaultValue(ctx context.Context, dftValue string) (interface{}, error)
}
