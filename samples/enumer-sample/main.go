package main

import (
	"encoding/json"
	"enumersample/enum"
	"fmt"

	"github.com/lopolopen/shoot"
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

	var x2 X
	json.Unmarshal(xJson, &x2)
	fmt.Println(x2)

	l, _ := shoot.ParseEnum[enum.Light]("Red")
	fmt.Println(l)

	if !shoot.IsEnum[enum.Light](5) {
		fmt.Println("not")
	}

	if shoot.TryParseEnum("Green", &l) {
		fmt.Println(l)
	}
}
