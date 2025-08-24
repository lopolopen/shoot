# shoot        [中文](https://github.com/lopolopen/shoot/blob/main/README.zh-cn.md)
A pack of tools for "go generate".

# Project status:
BETA VERSION. USE WITH CAUTION!

# How to use?

## using go1.23+ (SUGGESTED)

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

## using go1.23-

```zsh
# installing this tool needs go1.23+
# can be used by go1.23- project after installed
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

# TODO:
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

# Inspiring projects
* [stringer](https://pkg.go.dev/golang.org/x/tools/cmd/stringer)
* [enumer](https://github.com/dmarkham/enumer)
* [genapi](https://github.com/lexcao/genapi)
* [Refit](https://github.com/reactiveui/refit)