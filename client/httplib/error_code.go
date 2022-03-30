// Copyright 2020 beego
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

package httplib

import (
	"github.com/beego/beego/v2/core/berror"
)

var InvalidUrl = berror.DefineCode(4001001, moduleName, "InvalidUrl", `
You pass a invalid url to httplib module. Please check your url, be careful about special character. 
`)

var InvalidUrlProtocolVersion = berror.DefineCode(4001002, moduleName, "InvalidUrlProtocolVersion", `
You pass a invalid protocol version. In practice, we use HTTP/1.0, HTTP/1.1, HTTP/1.2
But something like HTTP/3.2 is valid for client, and the major version is 3, minor version is 2.
but you must confirm that server support those abnormal protocol version.
`)

var UnsupportedBodyType = berror.DefineCode(4001003, moduleName, "UnsupportedBodyType", `
You use a invalid data as request body.
For now, we only support type string and byte[].
`)

var InvalidXMLBody = berror.DefineCode(4001004, moduleName, "InvalidXMLBody", `
You pass invalid data which could not be converted to XML documents. In general, if you pass structure, it works well.
Sometimes you got XML document and you want to make it as request body. So you call XMLBody.
If you do this, you got this code. Instead, you should call Header to set Content-type and call Body to set body data.
`)

var InvalidYAMLBody = berror.DefineCode(4001005, moduleName, "InvalidYAMLBody", `
You pass invalid data which could not be converted to YAML documents. In general, if you pass structure, it works well.
Sometimes you got YAML document and you want to make it as request body. So you call YAMLBody.
If you do this, you got this code. Instead, you should call Header to set Content-type and call Body to set body data.
`)

var InvalidJSONBody = berror.DefineCode(4001006, moduleName, "InvalidJSONBody", `
You pass invalid data which could not be converted to JSON documents. In general, if you pass structure, it works well.
Sometimes you got JSON document and you want to make it as request body. So you call JSONBody.
If you do this, you got this code. Instead, you should call Header to set Content-type and call Body to set body data.
`)

var InvalidURLOrMethod = berror.DefineCode(4001007, moduleName, "InvalidURLOrMethod", `
You pass invalid url or method to httplib module. Please check the url and method, be careful about special characters.
`)

// start with 5 --------------------------------------------------------------------------

var CreateFormFileFailed = berror.DefineCode(5001001, moduleName, "CreateFormFileFailed", `
In normal case than handling files with BeegoRequest, you should not see this error code.
Unexpected EOF, invalid characters, bad file descriptor may cause this error.
`)

var ReadFileFailed = berror.DefineCode(5001002, moduleName, "ReadFileFailed", `
There are several cases that cause this error:
1. file not found. Please check the file name;
2. file not found, but file name is correct. If you use relative file path, it's very possible for you to see this code.
make sure that this file is in correct directory which Beego looks for;
3. Beego don't have the privilege to read this file, please change file mode; 
`)

var CopyFileFailed = berror.DefineCode(5001003, moduleName, "CopyFileFailed", `
When we try to read file content and then copy it to another writer, and failed.
1. Unexpected EOF;
2. Bad file descriptor;
3. Write conflict;

Please check your file content, and confirm that file is not processed by other process (or by user manually).
`)

var CloseFileFailed = berror.DefineCode(5001004, moduleName, "CloseFileFailed", `
After handling files, Beego try to close file but failed. Usually it was caused by bad file descriptor.
`)

var SendRequestFailed = berror.DefineCode(5001005, moduleName, "SendRequestRetryExhausted", `
Beego send HTTP request, but it failed.
If you config retry times, it means that Beego had retried and failed.
When you got this error, there are vary kind of reason:
1. Network unstable and timeout. In this case, sometimes server has received the request.
2. Server error. Make sure that server works well.
3. The request is invalid, which means that you pass some invalid parameter.
`)

var ReadGzipBodyFailed = berror.DefineCode(5001006, moduleName, "BuildGzipReaderFailed", `
Beego parse gzip-encode body failed. Usually Beego got invalid response.
Please confirm that server returns gzip data.
`)

var CreateFileIfNotExistFailed = berror.DefineCode(5001007, moduleName, "CreateFileIfNotExist", `
Beego want to create file if not exist and failed. 
In most cases, it means that Beego doesn't have the privilege to create this file.
Please change file mode to ensure that Beego is able to create files on specific directory.
Or you can run Beego with higher authority.
In some cases, you pass invalid filename. Make sure that the file name is valid on your system.
`)

var UnmarshalJSONResponseToObjectFailed = berror.DefineCode(5001008, moduleName,
	"UnmarshalResponseToObjectFailed", `
Beego trying to unmarshal response's body to structure but failed.
Make sure that:
1. You pass valid structure pointer to the function;
2. The body is valid json document
`)

var UnmarshalXMLResponseToObjectFailed = berror.DefineCode(5001009, moduleName,
	"UnmarshalResponseToObjectFailed", `
Beego trying to unmarshal response's body to structure but failed.
Make sure that:
1. You pass valid structure pointer to the function;
2. The body is valid XML document
`)

var UnmarshalYAMLResponseToObjectFailed = berror.DefineCode(5001010, moduleName,
	"UnmarshalResponseToObjectFailed", `
Beego trying to unmarshal response's body to structure but failed.
Make sure that:
1. You pass valid structure pointer to the function;
2. The body is valid YAML document
`)

var UnmarshalResponseToObjectFailed = berror.DefineCode(5001011, moduleName,
	"UnmarshalResponseToObjectFailed", `
Beego trying to unmarshal response's body to structure but failed.
There are several cases that cause this error:
1. You pass valid structure pointer to the function;
2. The body is valid json, Yaml or XML document
`)
