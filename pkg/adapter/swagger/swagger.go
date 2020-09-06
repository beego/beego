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
//
// Swagger™ is a project used to describe and document RESTful APIs.
//
// The Swagger specification defines a set of files required to describe such an API. These files can then be used by the Swagger-UI project to display the API and Swagger-Codegen to generate clients in various languages. Additional utilities can also take advantage of the resulting files, such as testing tools.
// Now in version 2.0, Swagger is more enabling than ever. And it's 100% open source software.

// Package swagger struct definition
package swagger

import (
	"github.com/astaxie/beego/pkg/server/web/swagger"
)

// Swagger list the resource
type Swagger swagger.Swagger

// Information Provides metadata about the API. The metadata can be used by the clients if needed.
type Information swagger.Information

// Contact information for the exposed API.
type Contact swagger.Contact

// License information for the exposed API.
type License swagger.License

// Item Describes the operations available on a single path.
type Item swagger.Item

// Operation Describes a single API operation on a path.
type Operation swagger.Operation

// Parameter Describes a single operation parameter.
type Parameter swagger.Parameter

// ParameterItems A limited subset of JSON-Schema's items object. It is used by parameter definitions that are not located in "body".
// http://swagger.io/specification/#itemsObject
type ParameterItems swagger.ParameterItems

// Schema Object allows the definition of input and output data types.
type Schema swagger.Schema

// Propertie are taken from the JSON Schema definition but their definitions were adjusted to the Swagger Specification
type Propertie swagger.Propertie

// Response as they are returned from executing this operation.
type Response swagger.Response

// Security Allows the definition of a security scheme that can be used by the operations
type Security swagger.Security

// Tag Allows adding meta data to a single tag that is used by the Operation Object
type Tag swagger.Tag

// ExternalDocs include Additional external documentation
type ExternalDocs swagger.ExternalDocs
