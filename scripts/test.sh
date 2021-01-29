#!/bin/bash

docker-compose -f "$(pwd)/scripts/test_docker_compose.yaml" up -d

export ORM_DRIVER=mysql
export TZ=UTC
export ORM_SOURCE="beego:test@tcp(localhost:13306)/orm_test?charset=utf8"

# wait for services in images ready
sleep 5

go test "$(pwd)/..."

# clear all container
docker-compose -f "$(pwd)/scripts/test_docker_compose.yaml" down


