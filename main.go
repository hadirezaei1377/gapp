package main

import (
	"fmt"
	"gapp/entity"
	"gapp/repository/mysql"
)

func main() {

}

func testUserMysqlRepo() {
	mysqlRepo := mysql.New()

	createdUser, err := mysqlRepo.Register(entity.User{
		ID:          0,
		PhoneNumber: "0927343",
		Name:        "Hossein Nazari",
	})

	if err != nil {
		fmt.Println("register user", err)
	} else {
		fmt.Println("created user", createdUser)
	}

	isUnique, err := mysqlRepo.IsPhoneNumberUnique(createdUser.PhoneNumber + "23")
	if err != nil {
		fmt.Println("unique err", err)
	}

	fmt.Println("isUnique", isUnique)
}
