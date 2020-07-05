#!/bin/bash

docker-compose -f test_docker_compose.yaml up -d

export ORM_DRIVER=mysql
export TZ=UTC
export ORM_SOURCE="beego:test@tcp(localhost:13306)/orm_test?charset=utf8"

go test ./...

# clear all container
docker-compose -f test_docker_compose.yaml down


