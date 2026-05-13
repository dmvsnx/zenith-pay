package mocks

import "github.com/savanyv/zenith-pay/internal/model"

type UserRepo struct {
	CreateFn        func(user *model.User) error
	GetByUsernameFn func(username string) (*model.User, error)
}

func (m *UserRepo) Create(user *model.User) error {
	return m.CreateFn(user)
}

func (m *UserRepo) GetByUsername(username string) (*model.User, error) {
	return m.GetByUsernameFn(username)
}
