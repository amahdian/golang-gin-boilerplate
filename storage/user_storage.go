package storage

import (
	"github.com/amahdian/golang-gin-boilerplate/domain/model"
)

type UserStorage interface {
	CrudStorage[*model.User]

	FindByEmail(email string) (*model.User, error)
}
