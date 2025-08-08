package main

import (
	"constructor_sample/model"
	"encoding/json"
	"fmt"
)

func main() {
	user := model.NewUser("John Doe", 30)
	user.SetAge(31)
	user.Gender = "Male"
	j, _ := json.Marshal(user)
	println(string(j))

	var newUser model.User
	if err := json.Unmarshal(j, &newUser); err != nil {
		panic(err)
	}
	println(newUser.Name(), newUser.Age(), newUser.Gender)
	fmt.Printf("%+v\n", newUser)
}
