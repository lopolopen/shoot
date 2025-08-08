# shoot
A pack of tools for "go generate".

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