# shoot        [English](https://github.com/lopolopen/shoot/blob/main/README.md)
为 "go generate" 而打造的工具集。

# 项目状态:
BETA版本，谨慎使用！

# 如何使用?

## 使用 go1.23+ （推荐）

```zsh
go get -tool github.com/lopolopen/shoot@latest
```

```go
//go:generate go tool shoot new -getset -json -type=YourType

//go:generate go tool shoot new -getset -json -file $GOFILE
```

```zsh
go generate
go generate ./...
```

## 使用 go1.23-

```zsh
# 安装此工具需要 go1.23 或以上版本
# 安装后可用于 go1.23 以下版本的项目
go install github.com/lopolopen/shoot@latest
```

```go
//go:generate shoot new -getset -json -type=YourType

//go:generate shoot new -getset -json -file $GOFILE
```

```zsh
go generate
go generate ./...
```

# 待办：
- [x] shoot new -getset -type=YourType
- [x] shoot new: field instruction like get;default=1
- [x] shoot new -opt|option -type=YourType
- [x] shoot new -json -type=YourType
- [x] shoot new: embed struct
- [x] shoot new: external package
- [x] shoot new -file=YourFile
- [ ] shoot new: type instruction like ignore
- [x] shoot new: -separate
- [ ] shoot enum -str|string -type=YourEnum
- [ ] shoot enum -bit|bitwise -type=YourEnum
- [ ] shoot enum -json -type=YourEnum
- [ ] refactor: duplicated code

# 启发项目：
* [stringer](https://pkg.go.dev/golang.org/x/tools/cmd/stringer)
* [enumer](https://github.com/dmarkham/enumer)
* [genapi](https://github.com/lexcao/genapi)
* [Refit](https://github.com/reactiveui/refit)
