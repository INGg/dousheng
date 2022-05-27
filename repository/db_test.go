package repository

import (
	"fmt"
	"testing"
)

func TestDB(t *testing.T) {
	InitDb()
	//CreateUser("123", "345")
	fmt.Println(UserCount)

	//_, _ = CreateUser("lzd", "213344")
	//_, _ = CreateUser("bob", "123444")
	//_, _ = CreateUser("alice", "98765")

	//u, _ := FindUserById(1)
	//fmt.Println(u)

	//fmt.Println(UserCount)
}

func TestFindAllVideoByUid(t *testing.T) {
	InitDb()
	id := 1
	var video []Video
	err := FindAllVideoByUid(int64(id), &video)
	if err != nil {
		return
	}
	fmt.Println(len(video))

	for i, v := range video {
		fmt.Println(i, v)
	}
}
