@ cmd

docker-compose -f test_docker_compose.yaml up -d

SET ORM_DRIVER=mysql
SET TZ=UTC
SET ORM_SOURCE=beego:test@tcp(localhost:13306)/orm_test?charset=utf8

go test ./...

@ clear all container
docker-compose -f test_docker_compose.yaml down


