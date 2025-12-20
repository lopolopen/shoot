package main

import (
	"fmt"
	"mappersample/domain/model"
	"mappersample/infra/po"
)

func main() {
	// orderDto := dto.NewOrder("1", "99", enums.OrderStatusCompleted)
	// order := orderDto.ToModel()
	// fmt.Println("Domain model.Order: ", order)

	// orderDto2 := new(dto.Order).FromModel(order)
	// fmt.Println("dto.Order: ", orderDto2)

	// userDto := dto.NewUser("1", "Zhen Chen", "zhen@chen.com")
	// user := userDto.ToModel()
	// fmt.Println("Domain model.User: ", user)

	// userDto2 := &dto.User{}
	// userDto2.FromModel(user)
	// fmt.Println("dto.User: ", userDto2)

	m := &model.User{
		ID:        11,
		FirstName: "yao",
		Model1: &model.Model1{
			Model2: model.Model2{
				B: 1111111,
			},
		},
	}
	p := new(po.User).FromDomain(m)
	fmt.Println(*p.Model, p)
}
