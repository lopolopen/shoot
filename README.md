# shoot
A pack of tools for "go generate".

# Project status:
PROVE OF CONCEPT. DON'T USE NOW!

# How to use?
```zsh
go get -tool github.com/lopolopen/shoot@latest
```

```go
//go:generate go tool shoot new -type=YourType

//go:generate go tool shoot enum -bit -type=YourEnum
```

```zsh
go generate
```

# TODO:
- [ ] shoot new -getset -type=YourType
- [ ] shoot new -opt|option -type=YourType
- [ ] shoot new -json -type=YourType
- [ ] shoot enum -str|string -type=YourEnum
- [x] shoot enum -bit|bitwise -type=YourEnum
- [ ] shoot enum -json -type=YourEnum
- [ ] refactor: duplicated code
