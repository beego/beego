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

// swagger struct definition
package swagger

const SwaggerVersion = "1.2"

type ResourceListing struct {
	ApiVersion     string `json:"apiVersion"`
	SwaggerVersion string `json:"swaggerVersion"` // e.g 1.2
	// BasePath       string `json:"basePath"`  obsolete in 1.1
	Apis  []ApiRef   `json:"apis"`
	Infos Infomation `json:"info"`
}

type ApiRef struct {
	Path        string `json:"path"` // relative or absolute, must start with /
	Description string `json:"description"`
}

type Infomation struct {
	Title             string `json:"title,omitempty"`
	Description       string `json:"description,omitempty"`
	Contact           string `json:"contact,omitempty"`
	TermsOfServiceUrl string `json:"termsOfServiceUrl,omitempty"`
	License           string `json:"license,omitempty"`
	LicenseUrl        string `json:"licenseUrl,omitempty"`
}

// https://github.com/wordnik/swagger-core/blob/scala_2.10-1.3-RC3/schemas/api-declaration-schema.json
type ApiDeclaration struct {
	ApiVersion     string           `json:"apiVersion"`
	SwaggerVersion string           `json:"swaggerVersion"`
	BasePath       string           `json:"basePath"`
	ResourcePath   string           `json:"resourcePath"` // must start with /
	Consumes       []string         `json:"consumes,omitempty"`
	Produces       []string         `json:"produces,omitempty"`
	Apis           []Api            `json:"apis,omitempty"`
	Models         map[string]Model `json:"models,omitempty"`
}

type Api struct {
	Path        string      `json:"path"` // relative or absolute, must start with /
	Description string      `json:"description"`
	Operations  []Operation `json:"operations,omitempty"`
}

type Operation struct {
	HttpMethod string `json:"httpMethod"`
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

type Protocol struct {
}

type ResponseMessage struct {
	Code          int    `json:"code"`
	Message       string `json:"message"`
	ResponseModel string `json:"responseModel"`
}

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

type ErrorResponse struct {
	Code   int    `json:"code"`
	Reason string `json:"reason"`
}

type Model struct {
	Id         string                   `json:"id"`
	Required   []string                 `json:"required,omitempty"`
	Properties map[string]ModelProperty `json:"properties"`
}

type ModelProperty struct {
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Items       map[string]string `json:"items,omitempty"`
	Format      string            `json:"format"`
}

// https://github.com/wordnik/swagger-core/wiki/authorizations
type Authorization struct {
	LocalOAuth OAuth  `json:"local-oauth"`
	ApiKey     ApiKey `json:"apiKey"`
}

// https://github.com/wordnik/swagger-core/wiki/authorizations
type OAuth struct {
	Type       string               `json:"type"`   // e.g. oauth2
	Scopes     []string             `json:"scopes"` // e.g. PUBLIC
	GrantTypes map[string]GrantType `json:"grantTypes"`
}

// https://github.com/wordnik/swagger-core/wiki/authorizations
type GrantType struct {
	LoginEndpoint        Endpoint `json:"loginEndpoint"`
	TokenName            string   `json:"tokenName"` // e.g. access_code
	TokenRequestEndpoint Endpoint `json:"tokenRequestEndpoint"`
	TokenEndpoint        Endpoint `json:"tokenEndpoint"`
}

// https://github.com/wordnik/swagger-core/wiki/authorizations
type Endpoint struct {
	Url              string `json:"url"`
	ClientIdName     string `json:"clientIdName"`
	ClientSecretName string `json:"clientSecretName"`
	TokenName        string `json:"tokenName"`
}

// https://github.com/wordnik/swagger-core/wiki/authorizations
type ApiKey struct {
	Type   string `json:"type"`   // e.g. apiKey
	PassAs string `json:"passAs"` // e.g. header
}
