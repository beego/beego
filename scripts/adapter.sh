#!/bin/sh

# using pkg/adapter. Usually you want to migrate to V2 smoothly, you could running this script

find ./ -name '*.go' -type f -exec sed -i '' -e 's/github.com\/astaxie\/beego/github.com\/astaxie\/beego\/pkg\/adapter/g' {} \;
find ./ -name '*.go' -type f -exec sed -i '' -e 's/"github.com\/astaxie\/beego\/pkg\/adapter"/beego "github.com\/astaxie\/beego\/pkg\/adapter"/g' {} \;
