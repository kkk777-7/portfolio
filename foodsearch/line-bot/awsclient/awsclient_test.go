package awsclient

import (
	"fmt"
	"line-bot/search"
	"os"
	"testing"
)

func TestIsLineUser(t *testing.T) {
	os.Setenv("AWSREGION", "us-east-1")
	os.Setenv("DYNAMOENDPOINT", "http://localhost:8000")
	var user User
	client := NewClient()

	err := client.IsLineUser("users", "123456789", &user)
	if err != nil && user.UserId != "" {
		fmt.Println(err)
	}
	fmt.Println(user.Status)
	fmt.Println(user)

	user2 := User{UserId: "1", Status: "NG"}
	err = client.SetLineUser("users", &user2)
	if err != nil {
		fmt.Println(err)
	}

	err = client.UpdateLineUser("users", &user2, "Genre", "イタリアン")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(user2)

	shop := search.Shop{UserId: "11111", ShopId: "22222", Genre: "イタリアン"}
	err = client.SetShop("shops", &shop)
	if err != nil {
		fmt.Println(err)
	}

	var shops []search.Shop
	err = client.GetShop("shops", "11111", &shops)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(shops)

	err = client.DeleteShop("shops", "11111", "22222")
}
