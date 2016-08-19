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

// Swagger list the resource
type Swagger struct {
	SwaggerVersion      string              `json:"swagger,omitempty"`
	Infos               Information         `json:"info"`
	Host                string              `json:"host,omitempty"`
	BasePath            string              `json:"basePath,omitempty"`
	Schemes             []string            `json:"schemes,omitempty"`
	Consumes            []string            `json:"consumes,omitempty"`
	Produces            []string            `json:"produces,omitempty"`
	Paths               map[string]*Item    `json:"paths"`
	Definitions         map[string]Schema   `json:"definitions,omitempty"`
	SecurityDefinitions map[string]Scurity  `json:"securityDefinitions,omitempty"`
	Security            map[string][]string `json:"security,omitempty"`
	Tags                []Tag               `json:"tags,omitempty"`
	ExternalDocs        *ExternalDocs       `json:"externalDocs,omitempty"`
}

// Information Provides metadata about the API. The metadata can be used by the clients if needed.
type Information struct {
	Title          string `json:"title,omitempty"`
	Description    string `json:"description,omitempty"`
	Version        string `json:"version,omitempty"`
	TermsOfService string `json:"termsOfService,omitempty"`

	Contact Contact `json:"contact,omitempty"`
	License License `json:"license,omitempty"`
}

// Contact information for the exposed API.
type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	EMail string `json:"email,omitempty"`
}

// License information for the exposed API.
type License struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// Item Describes the operations available on a single path.
type Item struct {
	Ref     string     `json:"$ref,omitempty"`
	Get     *Operation `json:"get,omitempty"`
	Put     *Operation `json:"put,omitempty"`
	Post    *Operation `json:"post,omitempty"`
	Delete  *Operation `json:"delete,omitempty"`
	Options *Operation `json:"options,omitempty"`
	Head    *Operation `json:"head,omitempty"`
	Patch   *Operation `json:"patch,omitempty"`
}

// Operation Describes a single API operation on a path.
type Operation struct {
	Tags        []string            `json:"tags,omitempty"`
	Summary     string              `json:"summary,omitempty"`
	Description string              `json:"description,omitempty"`
	OperationID string              `json:"operationId,omitempty"`
	Consumes    []string            `json:"consumes,omitempty"`
	Produces    []string            `json:"produces,omitempty"`
	Schemes     []string            `json:"schemes,omitempty"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	Responses   map[string]Response `json:"responses,omitempty"`
	Deprecated  bool                `json:"deprecated,omitempty"`
}

// Parameter Describes a single operation parameter.
type Parameter struct {
	In          string  `json:"in,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Required    bool    `json:"required,omitempty"`
	Schema      *Schema `json:"schema,omitempty"`
	Type        string  `json:"type,omitempty"`
	Format      string  `json:"format,omitempty"`
}

// Schema Object allows the definition of input and output data types.
type Schema struct {
	Ref         string               `json:"$ref,omitempty"`
	Title       string               `json:"title,omitempty"`
	Format      string               `json:"format,omitempty"`
	Description string               `json:"description,omitempty"`
	Required    []string             `json:"required,omitempty"`
	Type        string               `json:"type,omitempty"`
	Items       *Propertie           `json:"items,omitempty"`
	Properties  map[string]Propertie `json:"properties,omitempty"`
}

// Propertie are taken from the JSON Schema definition but their definitions were adjusted to the Swagger Specification
type Propertie struct {
	Ref                  string               `json:"$ref,omitempty"`
	Title                string               `json:"title,omitempty"`
	Description          string               `json:"description,omitempty"`
	Default              string               `json:"default,omitempty"`
	Type                 string               `json:"type,omitempty"`
	Example              string               `json:"example,omitempty"`
	Required             []string             `json:"required,omitempty"`
	Format               string               `json:"format,omitempty"`
	ReadOnly             bool                 `json:"readOnly,omitempty"`
	Properties           map[string]Propertie `json:"properties,omitempty"`
	Items                *Propertie           `json:"items,omitempty"`
	AdditionalProperties *Propertie           `json:"additionalProperties,omitempty"`
}

// Response as they are returned from executing this operation.
type Response struct {
	Description string  `json:"description,omitempty"`
	Schema      *Schema `json:"schema,omitempty"`
	Ref         string  `json:"$ref,omitempty"`
}

// Scurity Allows the definition of a security scheme that can be used by the operations
type Scurity struct {
	Type             string            `json:"type,omitempty"` // Valid values are "basic", "apiKey" or "oauth2".
	Description      string            `json:"description,omitempty"`
	Name             string            `json:"name,omitempty"`
	In               string            `json:"in,omitempty"`   // Valid values are "query" or "header".
	Flow             string            `json:"flow,omitempty"` // Valid values are "implicit", "password", "application" or "accessCode".
	AuthorizationURL string            `json:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty"` // The available scopes for the OAuth2 security scheme.
}

// Tag Allows adding meta data to a single tag that is used by the Operation Object
type Tag struct {
	Name         string        `json:"name,omitempty"`
	Description  string        `json:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}

// ExternalDocs include Additional external documentation
type ExternalDocs struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}
