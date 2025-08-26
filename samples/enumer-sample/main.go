package main

import (
	"encoding/json"
	"enumersample/enum"
	"fmt"
)

type X struct {
	Name  string
	Light enum.Light
}

func main() {
	fmt.Printf("LightRed: %s\n", enum.LightRed)

	x := X{
		Name:  "x",
		Light: enum.LightGreen,
	}
	xJson, _ := json.Marshal(x)
	fmt.Println(string(xJson))
}
