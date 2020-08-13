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

// BeanMetadata, in other words, bean's config.
// it could be read from config file
type BeanMetadata struct {
	// Fields: field name => field metadata
	Fields map[string]*FieldMetadata
}

// FieldMetadata contains metadata
type FieldMetadata struct {
	// default value in string format
	DftValue string
}
