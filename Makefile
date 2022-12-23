# Copyright 2020
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


##@ Test

test-orm-mysql5: ## Run ORM unit tests on mysql5.
	docker-compose -f scripts/orm_docker_compose.yaml up -d
	export ORM_DRIVER=mysql
	export ORM_SOURCE="beego:test@tcp(localhost:13306)/orm_test?charset=utf8"
	go test -v github.com/beego/beego/v2/client/orm
	docker-compose -f scripts/orm_docker_compose.yaml down

test-orm-mysql8: ## Run ORM unit tests on mysql8.
	docker-compose -f scripts/orm_docker_compose.yaml up -d
	export ORM_DRIVER=mysql
	export ORM_SOURCE="beego:test@tcp(localhost:23306)/orm_test?charset=utf8"
	go test -v github.com/beego/beego/v2/client/orm
	docker-compose -f scripts/orm_docker_compose.yaml down

test-orm-pgsql: ## Run ORM unit tests on postgresql.
	docker-compose -f scripts/orm_docker_compose.yaml up -d
	export ORM_DRIVER=postgres
	export ORM_SOURCE="user=postgres password=postgres dbname=orm_test sslmode=disable"
	go test -v github.com/beego/beego/v2/client/orm
	docker-compose -f scripts/orm_docker_compose.yaml down

test-orm-tidb: ## Run ORM unit tests on tidb.
	docker-compose -f scripts/orm_docker_compose.yaml up -d
	export ORM_DRIVER=tidb
	export ORM_SOURCE="memory://test/test"
	go test -v github.com/beego/beego/v2/client/orm
	docker-compose -f scripts/orm_docker_compose.yaml down

.PHONY: test-orm-all
test-orm-all: test-orm-mysql5 test-orm-mysql8 test-orm-pgsql test-orm-tidb

.PHONY: fmt
fmt:
	goimports -local "github.com/beego/beego" -w .