# Contributing to beego

beego is an open source project.

It is the work of hundreds of contributors. We appreciate your help!

Here are instructions to get you started. They are probably not perfect, please let us know if anything feels wrong or
incomplete.

## Prepare environment

Firstly, install some tools. Execute those commands **outside** the project. Or those command will modify go.mod file.

```shell script
go get -u golang.org/x/tools/cmd/goimports

go get -u github.com/gordonklaus/ineffassign
```

Put those lines into your pre-commit githook script:

```shell script
goimports -w -format-only ./

ineffassign .

staticcheck -show-ignored -checks "-ST1017,-U1000,-ST1005,-S1034,-S1012,-SA4006,-SA6005,-SA1019,-SA1024" ./
```

## Prepare middleware

Beego uses many middlewares, including MySQL, Redis, SSDB and so on.

We provide docker compose file to start all middlewares.

You can run:

```shell script
docker-compose -f scripts/test_docker_compose.yaml up -d
```

Unit tests read addresses from environment, here is an example:

```shell script
export ORM_DRIVER=mysql
export ORM_SOURCE="beego:test@tcp(192.168.0.105:13306)/orm_test?charset=utf8"
export MEMCACHE_ADDR="192.168.0.105:11211"
export REDIS_ADDR="192.168.0.105:6379"
export SSDB_ADDR="192.168.0.105:8888"
```

## Contribution guidelines

### Pull requests

First, beego follow the gitflow. So please send you pull request to **develop** branch. We will close the pull
request to master branch.

By the way, please don't forget update the `CHANGELOG.md` before you send pull request. 
You can just add your pull request following 'developing' section in `CHANGELOG.md`. 
We'll release them in the next Beego version.

We are always happy to receive pull requests, and do our best to review them as fast as possible. Not sure if that typo
is worth a pull request? Do it! We will appreciate it.

Don't forget to rebase your commits!

If your pull request is not accepted on the first try, don't be discouraged! Sometimes we can make a mistake, please do
more explaining for us. We will appreciate it.

We're trying very hard to keep beego simple and fast. We don't want it to do everything for everybody. This means that
we might decide against incorporating a new feature. But we will give you some advice on how to do it in other way.

### Create issues

Any significant improvement should be documented as [a GitHub issue](https://github.com/beego/beego/v2/issues) before
anybody starts working on it.

Also when filing an issue, make sure to answer these five questions:

- What version of beego are you using (bee version)?
- What operating system and processor architecture are you using?
- What did you do?
- What did you expect to see?
- What did you see instead?

### but check existing issues and docs first!

Please take a moment to check that an issue doesn't already exist documenting your bug report or improvement proposal.
If it does, it never hurts to add a quick "+1" or "I have this problem too". This will help prioritize the most common
problems and requests.

Also, if you don't know how to use it. please make sure you have read through the docs in http://beego.vip/docs
