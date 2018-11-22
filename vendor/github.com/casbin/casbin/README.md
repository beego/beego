Casbin
====

[![Go Report Card](https://goreportcard.com/badge/github.com/casbin/casbin)](https://goreportcard.com/report/github.com/casbin/casbin)
[![Build Status](https://travis-ci.org/casbin/casbin.svg?branch=master)](https://travis-ci.org/casbin/casbin)
[![Coverage Status](https://coveralls.io/repos/github/casbin/casbin/badge.svg?branch=master)](https://coveralls.io/github/casbin/casbin?branch=master)
[![Godoc](https://godoc.org/github.com/casbin/casbin?status.svg)](https://godoc.org/github.com/casbin/casbin)
[![Release](https://img.shields.io/github/release/casbin/casbin.svg)](https://github.com/casbin/casbin/releases/latest)
[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/casbin/lobby)
[![Patreon](https://img.shields.io/badge/patreon-donate-yellow.svg)](http://www.patreon.com/yangluo)
[![Sourcegraph Badge](https://sourcegraph.com/github.com/casbin/casbin/-/badge.svg)](https://sourcegraph.com/github.com/casbin/casbin?badge)

**News**: still worry about how to write the correct Casbin policy? ``Casbin online editor`` is coming to help! Try it at: http://casbin.org/editor/

![casbin Logo](casbin-logo.png)

Casbin is a powerful and efficient open-source access control library for Golang projects. It provides support for enforcing authorization based on various [access control models](https://en.wikipedia.org/wiki/Computer_security_model).

## All the languages supported by Casbin:

[![golang](https://casbin.org/docs/assets/langs/golang.png)](https://github.com/casbin/casbin) | [![java](https://casbin.org/docs/assets/langs/java.png)](https://github.com/casbin/jcasbin) | [![nodejs](https://casbin.org/docs/assets/langs/nodejs.png)](https://github.com/casbin/node-casbin) | [![php](https://casbin.org/docs/assets/langs/php.png)](https://github.com/php-casbin/php-casbin)
----|----|----|----
[Casbin](https://github.com/casbin/casbin) | [jCasbin](https://github.com/casbin/jcasbin) | [node-Casbin](https://github.com/casbin/node-casbin) | [PHP-Casbin](https://github.com/php-casbin/php-casbin)
production-ready | production-ready | production-ready | production-ready

## Table of contents

- [Supported models](#supported-models)
- [How it works?](#how-it-works)
- [Features](#features)
- [Installation](#installation)
- [Documentation](#documentation)
- [Online editor](#online-editor)
- [Tutorials](#tutorials)
- [Get started](#get-started)
- [Policy management](#policy-management)
- [Policy persistence](#policy-persistence)
- [Policy consistence between multiple nodes](#policy-consistence-between-multiple-nodes)
- [Role manager](#role-manager)
- [Multi-threading](#multi-threading)
- [Benchmarks](#benchmarks)
- [Examples](#examples)
- [How to use Casbin as a service?](#how-to-use-casbin-as-a-service)
- [Our adopters](#our-adopters)

## Supported models

1. [**ACL (Access Control List)**](https://en.wikipedia.org/wiki/Access_control_list)
2. **ACL with [superuser](https://en.wikipedia.org/wiki/Superuser)**
3. **ACL without users**: especially useful for systems that don't have authentication or user log-ins.
3. **ACL without resources**: some scenarios may target for a type of resources instead of an individual resource by using permissions like ``write-article``, ``read-log``. It doesn't control the access to a specific article or log.
4. **[RBAC (Role-Based Access Control)](https://en.wikipedia.org/wiki/Role-based_access_control)**
5. **RBAC with resource roles**: both users and resources can have roles (or groups) at the same time.
6. **RBAC with domains/tenants**: users can have different role sets for different domains/tenants.
7. **[ABAC (Attribute-Based Access Control)](https://en.wikipedia.org/wiki/Attribute-Based_Access_Control)**: syntax sugar like ``resource.Owner`` can be used to get the attribute for a resource.
8. **[RESTful](https://en.wikipedia.org/wiki/Representational_state_transfer)**: supports paths like ``/res/*``, ``/res/:id`` and HTTP methods like ``GET``, ``POST``, ``PUT``, ``DELETE``.
9. **Deny-override**: both allow and deny authorizations are supported, deny overrides the allow.
10. **Priority**: the policy rules can be prioritized like firewall rules.

## How it works?

In Casbin, an access control model is abstracted into a CONF file based on the **PERM metamodel (Policy, Effect, Request, Matchers)**. So switching or upgrading the authorization mechanism for a project is just as simple as modifying a configuration. You can customize your own access control model by combining the available models. For example, you can get RBAC roles and ABAC attributes together inside one model and share one set of policy rules.

The most basic and simplest model in Casbin is ACL. ACL's model CONF is:

```ini
# Request definition
[request_definition]
r = sub, obj, act

# Policy definition
[policy_definition]
p = sub, obj, act

# Policy effect
[policy_effect]
e = some(where (p.eft == allow))

# Matchers
[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act

```

An example policy for ACL model is like:

```
p, alice, data1, read
p, bob, data2, write
```

It means:

- alice can read data1
- bob can write data2

We also support multi-line mode by appending '\\'  in the end:  

```ini
# Matchers
[matchers]
m = r.sub == p.sub && r.obj == p.obj \ 
  && r.act == p.act
```

Further more, if you are using ABAC,  you can try operator `in` like following in Casbin **golang** edition (jCasbin and Node-Casbin are not supported yet):

```ini
# Matchers
[matchers]
m = r.obj == p.obj && r.act == p.act || r.obj in ('data2', 'data3')
```

But you **SHOULD** make sure that the length of the array is **MORE** than **1**, otherwise there will cause it to panic.

For more operators, you may take a look at [govaluate](https://github.com/Knetic/govaluate)

## Features

What Casbin does:

1. enforce the policy in the classic ``{subject, object, action}`` form or a customized form as you defined, both allow and deny authorizations are supported.
2. handle the storage of the access control model and its policy.
3. manage the role-user mappings and role-role mappings (aka role hierarchy in RBAC).
4. support built-in superuser like ``root`` or ``administrator``. A superuser can do anything without explict permissions.
5. multiple built-in operators to support the rule matching. For example, ``keyMatch`` can map a resource key ``/foo/bar`` to the pattern ``/foo*``.

What Casbin does NOT do:

1. authentication (aka verify ``username`` and ``password`` when a user logs in)
2. manage the list of users or roles. I believe it's more convenient for the project itself to manage these entities. Users usually have their passwords, and Casbin is not designed as a password container. However, Casbin stores the user-role mapping for the RBAC scenario. 

## Installation

```
go get github.com/casbin/casbin
```

## Documentation

https://casbin.org/docs/en/overview

## Online editor

You can also use the online editor (http://casbin.org/editor/) to write your Casbin model and policy in your web browser. It provides functionality such as ``syntax highlighting`` and ``code completion``, just like an IDE for a programming language.

## Tutorials

https://casbin.org/docs/en/tutorials

## Get started

1. New a Casbin enforcer with a model file and a policy file:

    ```go
    e := casbin.NewEnforcer("path/to/model.conf", "path/to/policy.csv")
    ```

Note: you can also initialize an enforcer with policy in DB instead of file, see [Persistence](#persistence) section for details.

2. Add an enforcement hook into your code right before the access happens:

    ```go
    sub := "alice" // the user that wants to access a resource.
    obj := "data1" // the resource that is going to be accessed.
    act := "read" // the operation that the user performs on the resource.

    if e.Enforce(sub, obj, act) == true {
        // permit alice to read data1
    } else {
        // deny the request, show an error
    }
    ```

3. Besides the static policy file, Casbin also provides API for permission management at run-time. For example, You can get all the roles assigned to a user as below:

    ```go
    roles := e.GetRoles("alice")
    ```

See [Policy management APIs](#policy-management) for more usage.

4. Please refer to the ``_test.go`` files for more usage.

## Policy management

Casbin provides two sets of APIs to manage permissions:

- [Management API](https://github.com/casbin/casbin/blob/master/management_api.go): the primitive API that provides full support for Casbin policy management. See [here](https://github.com/casbin/casbin/blob/master/management_api_test.go) for examples.
- [RBAC API](https://github.com/casbin/casbin/blob/master/rbac_api.go): a more friendly API for RBAC. This API is a subset of Management API. The RBAC users could use this API to simplify the code. See [here](https://github.com/casbin/casbin/blob/master/rbac_api_test.go) for examples.

We also provide a web-based UI for model management and policy management:

![model editor](https://hsluoyz.github.io/casbin/ui_model_editor.png)

![policy editor](https://hsluoyz.github.io/casbin/ui_policy_editor.png)

## Policy persistence

In Casbin, the policy storage is implemented as an adapter (aka middleware for Casbin). To keep light-weight, we don't put adapter code in the main library (except the default file adapter). A complete list of Casbin adapters is provided as below. Any 3rd-party contribution on a new adapter is welcomed, please inform us and I will put it in this list:)

Adapter | Type | Author | Description
----|------|----|----
[File Adapter (built-in)](https://casbin.org/docs/en/policy-storage#file-adapter-built-in) | File | Casbin | Persistence for [.CSV (Comma-Separated Values)](https://en.wikipedia.org/wiki/Comma-separated_values) files
[Filtered File Adapter (built-in)](https://github.com/casbin/casbin#policy-enforcement-at-scale) | File | [@faceless-saint](https://github.com/faceless-saint) | Persistence for [.CSV (Comma-Separated Values)](https://en.wikipedia.org/wiki/Comma-separated_values) files with policy subset loading support
[Xorm Adapter](https://github.com/casbin/xorm-adapter) | ORM | Casbin | MySQL, PostgreSQL, TiDB, SQLite, SQL Server, Oracle are supported by [Xorm](https://github.com/go-xorm/xorm/)
[Gorm Adapter](https://github.com/casbin/gorm-adapter) | ORM | Casbin | MySQL, PostgreSQL, Sqlite3, SQL Server are supported by [Gorm](https://github.com/jinzhu/gorm/)
[Beego ORM Adapter](https://github.com/casbin/beego-orm-adapter) | ORM | Casbin | MySQL, PostgreSQL, Sqlite3 are supported by [Beego ORM](https://beego.me/docs/mvc/model/overview.md)
[MongoDB Adapter](https://github.com/casbin/mongodb-adapter) | NoSQL | Casbin | Persistence for [MongoDB](https://www.mongodb.com)
[Cassandra Adapter](https://github.com/casbin/cassandra-adapter) | NoSQL | Casbin | Persistence for [Apache Cassandra DB](http://cassandra.apache.org)
[Consul Adapter](https://github.com/ankitm123/consul-adapter) | KV store | [@ankitm123](https://github.com/ankitm123) | Persistence for [HashiCorp Consul](https://www.consul.io/)
[Redis Adapter](https://github.com/casbin/redis-adapter) | KV store | Casbin | Persistence for [Redis](https://redis.io/)
[Protobuf Adapter](https://github.com/casbin/protobuf-adapter) | Stream | Casbin | Persistence for [Google Protocol Buffers](https://developers.google.com/protocol-buffers/)
[JSON Adapter](https://github.com/casbin/json-adapter) | String | Casbin | Persistence for [JSON](https://www.json.org/)
[String Adapter](https://github.com/qiangmzsx/string-adapter) | String | [@qiangmzsx](https://github.com/qiangmzsx) | Persistence for String
[RQLite Adapter](https://github.com/edomosystems/rqlite-adapter) | SQL | [EDOMO Systems](https://github.com/edomosystems) | Persistence for [RQLite](https://github.com/rqlite/rqlite/)
[PostgreSQL Adapter](https://github.com/going/casbin-postgres-adapter) | SQL | [Going](https://github.com/going) | Persistence for [PostgreSQL](https://www.postgresql.org/)
[RethinkDB Adapter](https://github.com/adityapandey9/rethinkdb-adapter) | NoSQL | [@adityapandey9](https://github.com/adityapandey9) | Persistence for [RethinkDB](https://rethinkdb.com/)
[DynamoDB Adapter](https://github.com/HOOQTV/dynacasbin) | NoSQL | [HOOQ](https://github.com/HOOQTV) | Persistence for [Amazon DynamoDB](https://aws.amazon.com/dynamodb/)
[Minio/AWS S3 Adapter](https://github.com/Soluto/casbin-minio-adapter) | Object storage | [Soluto](https://github.com/Soluto) | Persistence for [Minio](https://github.com/minio/minio) and [Amazon S3](https://aws.amazon.com/s3/)
[Bolt Adapter](https://github.com/wirepair/bolt-adapter) | KV store | [@wirepair](https://github.com/wirepair) | Persistence for [Bolt](https://github.com/boltdb/bolt)

For details of adapters, please refer to the documentation: https://casbin.org/docs/en/policy-storage

## Policy enforcement at scale

Some adapters support filtered policy management. This means that the policy loaded by Casbin is a subset of the policy in storage based on a given filter. This allows for efficient policy enforcement in large, multi-tenant environments when parsing the entire policy becomes a performance bottleneck.

To use filtered policies with a supported adapter, simply call the `LoadFilteredPolicy` method. The valid format for the filter parameter depends on the adapter used. To prevent accidental data loss, the `SavePolicy` method is disabled when a filtered policy is loaded.

For example, the following code snippet uses the built-in filtered file adapter and the RBAC model with domains. In this case, the filter limits the policy to a single domain. Any policy lines for domains other than `"domain1"` are omitted from the loaded policy:

```go
import (
    "github.com/casbin/casbin"
)

enforcer := casbin.NewEnforcer()

adapter := fileadapter.NewFilteredAdapter("examples/rbac_with_domains_policy.csv")
enforcer.InitWithAdapter("examples/rbac_with_domains_model.conf", adapter)

filter := &fileadapter.Filter{
    P: []string{"", "domain1"},
    G: []string{"", "", "domain1"},
}
enforcer.LoadFilteredPolicy(filter)

// The loaded policy now only contains the entries pertaining to "domain1".
```

## Policy consistence between multiple nodes

We support to use distributed messaging systems like [etcd](https://github.com/coreos/etcd) to keep consistence between multiple Casbin enforcer instances. So our users can concurrently use multiple Casbin enforcers to handle large number of permission checking requests.

Similar to policy storage adapters, we don't put watcher code in the main library. Any support for a new messaging system should be implemented as a watcher. A complete list of Casbin watchers is provided as below. Any 3rd-party contribution on a new watcher is welcomed, please inform us and I will put it in this list:)

Watcher | Type | Author | Description
----|------|----|----
[Etcd Watcher](https://github.com/casbin/etcd-watcher) | KV store | Casbin | Watcher for [etcd](https://github.com/coreos/etcd)
[NATS Watcher](https://github.com/Soluto/casbin-nats-watcher) | Messaging system | [Soluto](https://github.com/Soluto) | Watcher for [NATS](https://nats.io/)
[ZooKeeper Watcher](https://github.com/grepsr/casbin-zk-watcher) | KV store  | [Grepsr](https://github.com/grepsr) | Watcher for [Apache ZooKeeper](https://zookeeper.apache.org/)
[Redis Watcher](https://github.com/billcobbler/casbin-redis-watcher) | KV store  | [@billcobbler](https://github.com/billcobbler) | Watcher for [Redis](http://redis.io/)

## Role manager

The role manager is used to manage the RBAC role hierarchy (user-role mapping) in Casbin. A role manager can retrieve the role data from Casbin policy rules or external sources such as LDAP, Okta, Auth0, Azure AD, etc. We support different implementations of a role manager. To keep light-weight, we don't put role manager code in the main library (except the default role manager). A complete list of Casbin role managers is provided as below. Any 3rd-party contribution on a new role manager is welcomed, please inform us and I will put it in this list:)

Role manager | Author | Description
----|----|----
[Default Role Manager (built-in)](https://github.com/casbin/casbin/blob/master/rbac/default-role-manager/role_manager.go) | Casbin | Supports role hierarchy stored in Casbin policy
[Session Role Manager](https://github.com/casbin/session-role-manager) | [EDOMO Systems](https://github.com/edomosystems) | Supports role hierarchy stored in Casbin policy, with time-range-based sessions
[Okta Role Manager](https://github.com/casbin/okta-role-manager) | Casbin | Supports role hierarchy stored in [Okta](https://www.okta.com/)
[Auth0 Role Manager](https://github.com/casbin/auth0-role-manager) | Casbin | Supports role hierarchy stored in [Auth0](https://auth0.com/)'s [Authorization Extension](https://auth0.com/docs/extensions/authorization-extension/v2)

For developers: all role managers must implement the [RoleManager](https://github.com/casbin/casbin/blob/master/rbac/role_manager.go) interface. [Session Role Manager](https://github.com/casbin/session-role-manager) can be used as a reference implementation.

## Multi-threading

If you use Casbin in a multi-threading manner, you can use the synchronized wrapper of the Casbin enforcer: https://github.com/casbin/casbin/blob/master/enforcer_synced.go.

It also supports the ``AutoLoad`` feature, which means the Casbin enforcer will automatically load the latest policy rules from DB if it has changed. Call ``StartAutoLoadPolicy()`` to start automatically loading policy periodically and call ``StopAutoLoadPolicy()`` to stop it.

## Benchmarks

The overhead of policy enforcement is benchmarked in [model_b_test.go](https://github.com/casbin/casbin/blob/master/model_b_test.go). The testbed is:

```
Intel(R) Core(TM) i7-6700HQ CPU @ 2.60GHz, 2601 Mhz, 4 Core(s), 8 Logical Processor(s)
```

The benchmarking result of ``go test -bench=. -benchmem`` is as follows (op = an ``Enforce()`` call, ms = millisecond, KB = kilo bytes):

Test case | Size | Time overhead | Memory overhead
----|------|------|----
ACL | 2 rules (2 users) | 0.015493 ms/op | 5.649 KB
RBAC | 5 rules (2 users, 1 role) | 0.021738 ms/op | 7.522 KB
RBAC (small) | 1100 rules (1000 users, 100 roles) | 0.164309 ms/op | 80.620 KB
RBAC (medium) | 11000 rules (10000 users, 1000 roles) | 2.258262 ms/op | 765.152 KB
RBAC (large) | 110000 rules (100000 users, 10000 roles) | 23.916776 ms/op | 7.606 MB
RBAC with resource roles | 6 rules (2 users, 2 roles) | 0.021146 ms/op | 7.906 KB
RBAC with domains/tenants | 6 rules (2 users, 1 role, 2 domains) | 0.032696 ms/op | 10.755 KB
ABAC | 0 rule (0 user) | 0.007510 ms/op | 2.328 KB
RESTful | 5 rules (3 users) | 0.045398 ms/op | 91.774 KB
Deny-override | 6 rules (2 users, 1 role) | 0.023281 ms/op | 8.370 KB
Priority | 9 rules (2 users, 2 roles) | 0.016389 ms/op | 5.313 KB

## Examples

Model | Model file | Policy file
----|------|----
ACL | [basic_model.conf](https://github.com/casbin/casbin/blob/master/examples/basic_model.conf) | [basic_policy.csv](https://github.com/casbin/casbin/blob/master/examples/basic_policy.csv)
ACL with superuser | [basic_model_with_root.conf](https://github.com/casbin/casbin/blob/master/examples/basic_with_root_model.conf) | [basic_policy.csv](https://github.com/casbin/casbin/blob/master/examples/basic_policy.csv)
ACL without users | [basic_model_without_users.conf](https://github.com/casbin/casbin/blob/master/examples/basic_without_users_model.conf) | [basic_policy_without_users.csv](https://github.com/casbin/casbin/blob/master/examples/basic_without_users_policy.csv)
ACL without resources | [basic_model_without_resources.conf](https://github.com/casbin/casbin/blob/master/examples/basic_without_resources_model.conf) | [basic_policy_without_resources.csv](https://github.com/casbin/casbin/blob/master/examples/basic_without_resources_policy.csv)
RBAC | [rbac_model.conf](https://github.com/casbin/casbin/blob/master/examples/rbac_model.conf)  | [rbac_policy.csv](https://github.com/casbin/casbin/blob/master/examples/rbac_policy.csv)
RBAC with resource roles | [rbac_model_with_resource_roles.conf](https://github.com/casbin/casbin/blob/master/examples/rbac_with_resource_roles_model.conf)  | [rbac_policy_with_resource_roles.csv](https://github.com/casbin/casbin/blob/master/examples/rbac_with_resource_roles_policy.csv)
RBAC with domains/tenants | [rbac_model_with_domains.conf](https://github.com/casbin/casbin/blob/master/examples/rbac_with_domains_model.conf)  | [rbac_policy_with_domains.csv](https://github.com/casbin/casbin/blob/master/examples/rbac_with_domains_policy.csv)
ABAC | [abac_model.conf](https://github.com/casbin/casbin/blob/master/examples/abac_model.conf)  | N/A
RESTful | [keymatch_model.conf](https://github.com/casbin/casbin/blob/master/examples/keymatch_model.conf)  | [keymatch_policy.csv](https://github.com/casbin/casbin/blob/master/examples/keymatch_policy.csv)
Deny-override | [rbac_model_with_deny.conf](https://github.com/casbin/casbin/blob/master/examples/rbac_with_deny_model.conf)  | [rbac_policy_with_deny.csv](https://github.com/casbin/casbin/blob/master/examples/rbac_with_deny_policy.csv)
Priority | [priority_model.conf](https://github.com/casbin/casbin/blob/master/examples/priority_model.conf)  | [priority_policy.csv](https://github.com/casbin/casbin/blob/master/examples/priority_policy.csv)

## How to use Casbin as a service?

- [Casbin Server](https://github.com/casbin/casbin-server): The official ``Casbin as a Service`` solution based on [gRPC](https://grpc.io/), both Management API and RBAC API are provided.
- [Go-Simple-API-Gateway](https://github.com/Soontao/go-simple-api-gateway): A simple API gateway written by golang, supports for authentication and authorization.
- [middleware-acl](https://github.com/luk4z7/middleware-acl): RESTful access control middleware based on Casbin.

## Our adopters

### Web frameworks

- [Beego](https://github.com/astaxie/beego): An open-source, high-performance web framework for Go, via built-in plugin: [plugins/authz](https://github.com/astaxie/beego/blob/master/plugins/authz)
- [Caddy](https://github.com/mholt/caddy): Fast, cross-platform HTTP/2 web server with automatic HTTPS, via plugin: [caddy-authz](https://github.com/casbin/caddy-authz)
- [Gin](https://github.com/gin-gonic/gin): A HTTP web framework featuring a Martini-like API with much better performance, via plugin: [authz](https://github.com/gin-contrib/authz)
- [Revel](https://github.com/revel/revel): A high productivity, full-stack web framework for the Go language, via plugin: [auth/casbin](https://github.com/revel/modules/tree/master/auth/casbin)
- [Echo](https://github.com/labstack/echo): High performance, minimalist Go web framework, via plugin: [echo-authz](https://github.com/labstack/echo-contrib/tree/master/casbin) (thanks to [@xqbumu](https://github.com/xqbumu))
- [Iris](https://github.com/kataras/iris): The fastest web framework for Go in (THIS) Earth. HTTP/2 Ready-To-GO, via plugin: [casbin](https://github.com/iris-contrib/middleware/tree/master/casbin) (thanks to [@hiveminded](https://github.com/hiveminded))
- [Negroni](https://github.com/urfave/negroni): Idiomatic HTTP Middleware for Golang, via plugin: [negroni-authz](https://github.com/casbin/negroni-authz)
- [Tango](https://github.com/lunny/tango): Micro & pluggable web framework for Go, via plugin: [authz](https://github.com/tango-contrib/authz)
- [Chi](https://github.com/pressly/chi): A lightweight, idiomatic and composable router for building HTTP services, via plugin: [chi-authz](https://github.com/casbin/chi-authz)
- [Macaron](https://github.com/go-macaron/macaron): A high productive and modular web framework in Go, via plugin: [authz](https://github.com/go-macaron/authz)
- [DotWeb](https://github.com/devfeel/dotweb): Simple and easy go web micro framework, via plugin: [authz](https://github.com/devfeel/middleware/tree/master/authz)
- [Baa](https://github.com/go-baa/baa): An express Go web framework with routing, middleware, dependency injection and http context, via plugin: [authz](https://github.com/baa-middleware/authz)

### Others

- [Intel RMD](https://github.com/intel/rmd): Intel's resource management daemon, via direct integration, see: [model](https://github.com/intel/rmd/blob/master/etc/rmd/acl/url/model.conf), [policy rules](https://github.com/intel/rmd/blob/master/etc/rmd/acl/url/policy.csv)
- [VMware Dispatch](https://github.com/vmware/dispatch): A framework for deploying and managing serverless style applications, via direct integration, see: [model (in code)](https://github.com/vmware/dispatch/blob/master/pkg/identity-manager/handlers.go#L46-L55), [policy rules (in code)](https://github.com/vmware/dispatch/blob/master/pkg/identity-manager/handlers_test.go#L35-L45)
- [Banzai Pipeline](https://github.com/banzaicloud/pipeline): [Banzai Cloud](https://github.com/banzaicloud)'s RESTful API to provision or reuse managed Kubernetes clusters in the cloud, via direct integration, see: [model (in code)](https://github.com/banzaicloud/pipeline/blob/master/auth/authz.go#L15-L30), [policy rules (in code)](https://github.com/banzaicloud/pipeline/blob/master/auth/authz.go#L84-L93)
- [Docker](https://github.com/docker/docker): The world's leading software container platform, via plugin: [casbin-authz-plugin](https://github.com/casbin/casbin-authz-plugin) ([recommended by Docker](https://docs.docker.com/engine/extend/legacy_plugins/#authorization-plugins))
- [Gobis](https://github.com/orange-cloudfoundry/gobis): [Orange](https://github.com/orange-cloudfoundry)'s lightweight API Gateway written in go, via plugin: [casbin](https://github.com/orange-cloudfoundry/gobis-middlewares/tree/master/casbin), see [model (in code)](https://github.com/orange-cloudfoundry/gobis-middlewares/blob/master/casbin/model.go#L52-L65), [policy rules (from request)](https://github.com/orange-cloudfoundry/gobis-middlewares/blob/master/casbin/adapter.go#L46-L64)
- [Skydive](https://github.com/skydive-project/skydive): An open source real-time network topology and protocols analyzer, via direct integration, see: [model (in code)](https://github.com/skydive-project/skydive/blob/master/config/config.go#L136-L140), [policy rules](https://github.com/skydive-project/skydive/blob/master/rbac/policy.csv)
- [Zenpress](https://github.com/insionng/zenpress): A CMS system written in Golang, via direct integration, see: [model](https://github.com/insionng/zenpress/blob/master/content/config/rbac_model.conf), [policy rules (in Gorm)](https://github.com/insionng/zenpress/blob/master/model/user.go#L53-L77)
- [Argo CD](https://github.com/argoproj/argo-cd): GitOps continuous delivery for Kubernetes, via direct integration, see: [model](https://github.com/argoproj/argo-cd/blob/master/util/rbac/model.conf), [policy rules](https://github.com/argoproj/argo-cd/blob/master/util/rbac/builtin-policy.csv)
- [Muxi Cloud](https://github.com/muxiyun/Mae): PaaS of Muxi Cloud, an easier way to manage Kubernetes cluster, via direct integration, see: [model](https://github.com/muxiyun/Mae/blob/master/conf/casbinmodel.conf), [policy rules (in code)](https://github.com/muxiyun/Mae/blob/master/pkg/casbin/initPolicy.go#L21-L95)
- [EngineerCMS](https://github.com/3xxx/EngineerCMS): A CMS to manage knowledge for engineers, via direct integration, see: [model](https://github.com/3xxx/EngineerCMS/blob/master/conf/rbac_model.conf), [policy rules (in SQLite)](https://github.com/3xxx/EngineerCMS/blob/master/database/engineer.db)
- [Cyber Auth API](https://github.com/CyberlifeCN/cyber-auth-api): A Golang authentication API project, via direct integration, see: [model](https://github.com/CyberlifeCN/cyber-auth-api/blob/master/conf/authz_model.conf), [policy rules](https://github.com/CyberlifeCN/cyber-auth-api/blob/master/conf/authz_policy.csv)
- [IRIS Community](https://github.com/irisnet/iris-community): Website for IRIS Community Activities, via direct integration, see: [model](https://github.com/irisnet/iris-community/blob/master/authz/authz_model.conf), [policy rules](https://github.com/irisnet/iris-community/blob/master/authz/authz_policy.csv)
- [Metadata DB](https://github.com/Bnei-Baruch/mdb): BB archive metadata database, via direct integration, see: [model](https://github.com/Bnei-Baruch/mdb/blob/master/data/permissions_model.conf), [policy rules](https://github.com/Bnei-Baruch/mdb/blob/master/data/permissions_policy.csv)

## License

This project is licensed under the [Apache 2.0 license](LICENSE).

## Contact

If you have any issues or feature requests, please contact us. PR is welcomed.
- https://github.com/casbin/casbin/issues
- hsluoyz@gmail.com
- Tencent QQ group: [546057381](//shang.qq.com/wpa/qunwpa?idkey=8ac8b91fc97ace3d383d0035f7aa06f7d670fd8e8d4837347354a31c18fac885)
