package main

import (
	"constructorsample/model"
	"encoding/json"
	"fmt"
	"regexp"
)

func main() {
	doc := "shoot: default=\"80\";"
	regDef := regexp.MustCompile(`(?im)^shoot:.*?\Wdef(ault)?=([^;]+)(;.*|\s*)$`)
	ms := regDef.FindStringSubmatch(doc)
	for _, m := range ms {
		fmt.Println(m)
	}

	user := model.NewUser("11", "Tom", 0, 21, "123456")
	fmt.Printf("%+v\n", user)

	user.SetAge(31)
	user.SetTel("999999")

	userJson, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(userJson))

	var user2 model.User
	err = json.Unmarshal(userJson, &user2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", &user2)

	conf := model.NewConf("1", "127.0.0.1", 80)
	fmt.Printf("%+v\n", conf)

	conf2 := model.NewConf2("2").With()
	fmt.Printf("%+v\n", conf2)

	//only init with default values
	conf2 = new(model.Conf2).With()
	fmt.Printf("%+v\n", conf2)
}
