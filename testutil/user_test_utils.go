package testutil

import (
	"github.com/amahdian/golang-gin-boilerplate/domain/model"
	"github.com/amahdian/golang-gin-boilerplate/global/test"
)

func TestUser() *model.User {
	return &model.User{
		Email: test.UserEmail,
	}
}
