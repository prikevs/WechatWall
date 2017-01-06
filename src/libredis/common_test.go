package libredis

import (
	"testing"
)

func TestUserJson(t *testing.T) {
	user := &User{
		UserOpenid:     "123",
		UserName:       "kevince",
		UserCreateTime: 9987,
	}
	js, err := user.Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(js)
}

func TestUserLoads(t *testing.T) {
	js := `{"user_openid":"123","user_name":"kevince","user_remark":"","user_create_time":9987}`
	user := &User{}
	if err := user.Loads(js); err != nil {
		t.Fatal(err)
	}
	t.Log(user)
}

func TestSetUserToMap(t *testing.T) {
	user := &User{
		UserOpenid:     "123",
		UserName:       "kevince",
		UserCreateTime: 9987,
	}
	mp, err := GetMap("test:map:1:")
	if err != nil {
		t.Fatal(err)
	}
	if err := SetUserToMap(user, mp); err != nil {
		t.Fatal(err)
	}
}
