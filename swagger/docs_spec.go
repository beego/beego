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

// Package swagger struct definition
package swagger

// SwaggerVersion show the current swagger version
const SwaggerVersion = "1.2"

// ResourceListing list the resource
type ResourceListing struct {
	APIVersion     string `json:"apiVersion"`
	SwaggerVersion string `json:"swaggerVersion"` // e.g 1.2
	// BasePath       string `json:"basePath"`  obsolete in 1.1
	APIs []APIRef    `json:"apis"`
	Info Information `json:"info"`
}

// APIRef description the api path and description
type APIRef struct {
	Path        string `json:"path"` // relative or absolute, must start with /
	Description string `json:"description"`
}

// Information show the API Information
type Information struct {
	Title             string `json:"title,omitempty"`
	Description       string `json:"description,omitempty"`
	Contact           string `json:"contact,omitempty"`
	TermsOfServiceURL string `json:"termsOfServiceUrl,omitempty"`
	License           string `json:"license,omitempty"`
	LicenseURL        string `json:"licenseUrl,omitempty"`
}

// APIDeclaration see https://github.com/wordnik/swagger-core/blob/scala_2.10-1.3-RC3/schemas/api-declaration-schema.json
type APIDeclaration struct {
	APIVersion     string           `json:"apiVersion"`
	SwaggerVersion string           `json:"swaggerVersion"`
	BasePath       string           `json:"basePath"`
	ResourcePath   string           `json:"resourcePath"` // must start with /
	Consumes       []string         `json:"consumes,omitempty"`
	Produces       []string         `json:"produces,omitempty"`
	APIs           []API            `json:"apis,omitempty"`
	Models         map[string]Model `json:"models,omitempty"`
}

// API show tha API struct
type API struct {
	Path        string      `json:"path"` // relative or absolute, must start with /
	Description string      `json:"description"`
	Operations  []Operation `json:"operations,omitempty"`
}

// Operation desc the Operation
type Operation struct {
	HTTPMethod string `json:"httpMethod"`
	Nickname   string `json:"nickname"`
	Type       string `json:"type"` // in 1.1 = DataType
	// ResponseClass    string            `json:"responseClass"` obsolete in 1.2
	Summary          string            `json:"summary,omitempty"`
	Notes            string            `json:"notes,omitempty"`
	Parameters       []Parameter       `json:"parameters,omitempty"`
	ResponseMessages []ResponseMessage `json:"responseMessages,omitempty"` // optional
	Consumes         []string          `json:"consumes,omitempty"`
	Produces         []string          `json:"produces,omitempty"`
	Authorizations   []Authorization   `json:"authorizations,omitempty"`
	Protocols        []Protocol        `json:"protocols,omitempty"`
}

// Protocol support which Protocol
type Protocol struct {
}

// ResponseMessage Show the
type ResponseMessage struct {
	Code          int    `json:"code"`
	Message       string `json:"message"`
	ResponseModel string `json:"responseModel"`
}

// Parameter desc the request parameters
type Parameter struct {
	ParamType     string `json:"paramType"` // path,query,body,header,form
	Name          string `json:"name"`
	Description   string `json:"description"`
	DataType      string `json:"dataType"` // 1.2 needed?
	Type          string `json:"type"`     // integer
	Format        string `json:"format"`   // int64
	AllowMultiple bool   `json:"allowMultiple"`
	Required      bool   `json:"required"`
	Minimum       int    `json:"minimum"`
	Maximum       int    `json:"maximum"`
}

// ErrorResponse desc response
type ErrorResponse struct {
	Code   int    `json:"code"`
	Reason string `json:"reason"`
}

// Model define the data model
type Model struct {
	ID         string                   `json:"id"`
	Required   []string                 `json:"required,omitempty"`
	Properties map[string]ModelProperty `json:"properties"`
}

// ModelProperty define the properties
type ModelProperty struct {
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Items       map[string]string `json:"items,omitempty"`
	Format      string            `json:"format"`
}

// Authorization see https://github.com/wordnik/swagger-core/wiki/authorizations
type Authorization struct {
	LocalOAuth OAuth  `json:"local-oauth"`
	APIKey     APIKey `json:"apiKey"`
}

// OAuth see https://github.com/wordnik/swagger-core/wiki/authorizations
type OAuth struct {
	Type       string               `json:"type"`   // e.g. oauth2
	Scopes     []string             `json:"scopes"` // e.g. PUBLIC
	GrantTypes map[string]GrantType `json:"grantTypes"`
}

// GrantType see https://github.com/wordnik/swagger-core/wiki/authorizations
type GrantType struct {
	LoginEndpoint        Endpoint `json:"loginEndpoint"`
	TokenName            string   `json:"tokenName"` // e.g. access_code
	TokenRequestEndpoint Endpoint `json:"tokenRequestEndpoint"`
	TokenEndpoint        Endpoint `json:"tokenEndpoint"`
}

// Endpoint see https://github.com/wordnik/swagger-core/wiki/authorizations
type Endpoint struct {
	URL              string `json:"url"`
	ClientIDName     string `json:"clientIdName"`
	ClientSecretName string `json:"clientSecretName"`
	TokenName        string `json:"tokenName"`
}

// APIKey see https://github.com/wordnik/swagger-core/wiki/authorizations
type APIKey struct {
	Type   string `json:"type"`   // e.g. apiKey
	PassAs string `json:"passAs"` // e.g. header
}
