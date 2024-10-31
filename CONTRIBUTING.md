# Contributing to beego

Beego is an open-source project.

It is the work of hundreds of contributors. And you could be among them, so we appreciate your help!

Here are instructions to get you started. They are probably not perfect, so please let us know if anything feels wrong or
incomplete.

## Prepare environment

Firstly, you need to install some tools. Execute the commands below **outside** the project. Otherwise, this action will modify the go.mod file.

```shell script
go get -u golang.org/x/tools/cmd/goimports

go get -u github.com/gordonklaus/ineffassign
```

Put the lines below in your pre-commit git hook script:

```shell script
goimports -w -format-only ./

ineffassign .

staticcheck -show-ignored -checks "-ST1017,-U1000,-ST1005,-S1034,-S1012,-SA4006,-SA6005,-SA1019,-SA1024" ./
```

## Prepare middleware

Beego uses many middlewares, including MySQL, Redis, SSDB amongs't others.

We provide a docker-compose file to start all middlewares.

You can run the following command to start all middlewares:

```shell script
docker-compose -f scripts/test_docker_compose.yaml up -d
```

Unit tests read addresses from environmental variables, you can set them up as shown in the example below:

```shell script
export ORM_DRIVER=mysql
export ORM_SOURCE="beego:test@tcp(192.168.0.105:13306)/orm_test?charset=utf8"
export MEMCACHE_ADDR="192.168.0.105:11211"
export REDIS_ADDR="192.168.0.105:6379"
export SSDB_ADDR="192.168.0.105:8888"
```

## Contribution guidelines

### Pull requests

Beego follows the gitflow. And as such, please submit your pull request to the **develop** branch. We will close the pull request by merging it into the master branch.

**NOTE:** Don't forget to update the `CHANGELOG.md` file by adding the changes made under the **developing** section.
We'll release them in the next Beego version.

We are always happy to receive pull requests, and do our best to review them as fast as possible. Not sure if that typo is worth a pull request? Just do it! We will appreciate it.

Don't forget to rebase your commits!

If your pull request is rejected, dont be discouraged. Sometimes we make mistakes. You can provide us with more context by explaining your issue as clearly as possible.

In our pursuit of maintaining Beego's simplicity and speed, we might not accept some feature requests. We don't want it to do everything for everybody. For this reason, we might decide against incorporating a new feature. However, we will provide guidance on achieving the same thing using a different approach

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

Take a moment to check that an issue documenting your bug report or improvement proposal doesn't already exist.
If it does, it doesn't hurts to add a quick "+1" or "I have this problem too". This will help prioritize the most common
problems and requests.
