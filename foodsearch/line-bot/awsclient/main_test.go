package awsclient

import (
	"fmt"
	"os"
	"testing"
)

func TestIsLineUser(t *testing.T) {
	os.Setenv("AWSREGION", "us-east-1")
	os.Setenv("DYNAMOENDPOINT", "http://localhost:8000")
	var user User
	client := NewClient("users")

	err := client.IsLineUser("123456789", &user)
	if err != nil && user.UserId != "" {
		fmt.Println(err)
	}
	fmt.Println(user.Status)
	fmt.Println(user)

	user2 := User{UserId: "1", Status: "NG"}
	err = client.SetLineUser(&user2)
	if err != nil {
		fmt.Println(err)
	}

	err = client.UpdateLineUser(&user2, "Genre", "イタリアン")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(user2)
}
