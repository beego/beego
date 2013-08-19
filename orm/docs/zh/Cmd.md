## 命令模式

注册模型与数据库以后，调用 RunCommand 执行 orm 命令

```go
func main() {
	// orm.RegisterModel...
	// orm.RegisterDataBase...
	...
	orm.RunCommand()
}
```

```bash
go build main.go
./main orm
# 直接执行可以显示帮助
# 如果你的程序可以支持的话，直接运行 go run main.go orm 也是一样的效果
```

## 自动建表

```bash
./main orm syncdb -h
Usage of orm command: syncdb:
  -db="default": DataBase alias name
  -force=false: drop tables before create
  -v=false: verbose info
```

使用 `-force=1` 可以 drop table 后再建表

使用 `-v` 可以查看执行的 sql 语句

## 打印建表SQL

```bash
./main orm sqlall -h
Usage of orm command: syncdb:
  -db="default": DataBase alias name
```

默认使用别名为 default 的数据库
