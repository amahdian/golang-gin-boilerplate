package svc

import (
	"context"

	"github.com/amahdian/golang-gin-boilerplate/global/env"

	"github.com/amahdian/golang-gin-boilerplate/storage"
)

type Svc interface {
	NewUserSvc(ctx context.Context) UserSvc
}

type svcImpl struct {
	stg  storage.Storage
	Envs *env.Envs
}

func NewSvc(stg storage.Storage, envs *env.Envs) Svc {
	return &svcImpl{
		stg,
		envs,
	}
}

func (s *svcImpl) NewUserSvc(ctx context.Context) UserSvc {
	return newUserSvc(ctx, s.stg, s.Envs)
}
