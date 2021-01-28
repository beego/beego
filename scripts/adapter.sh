#!/bin/sh

# using pkg/adapter. Usually you want to migrate to V2 smoothly, you could running this script

find ./ -name '*.go' -type f -exec sed -i '' -e 's/github.com\/astaxie\/beego/github.com\/astaxie\/beego\/pkg\/adapter/g' {} \;
find ./ -name '*.go' -type f -exec sed -i '' -e 's/"github.com\/astaxie\/beego\/pkg\/adapter"/beego "github.com\/astaxie\/beego\/pkg\/adapter"/g' {} \;

update rrp_flow set status = 4 where flow_id in (5623711176,5629411891)

select * from rrp_flow where flow_id in (5623711176,5629411891) and status = 5

update rrp_flow set status = 5 where flow_id in (5623711176,5629411891)

