package usecase

import (
	"errors"

	dtos "github.com/savanyv/zenith-pay/internal/dto"
	"github.com/savanyv/zenith-pay/internal/model"
	"github.com/savanyv/zenith-pay/internal/repository"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
)

type UserUsecase interface {
	Register(req *dtos.CreateUserRequest) (*dtos.CreateUserResponse, error)
	Login(req *dtos.LoginRequest) (*dtos.LoginResponse, error)
	ListUsers(page, limit int) ([]*dtos.UserResponse, int64, error)
}

type userUsecase struct {
	userRepo repository.UserRespository
	bcrypt helpers.BcryptHelper
	jwt helpers.JWTService
}

func NewUserUsecase(userRepo repository.UserRespository, jwt helpers.JWTService, bcrypt helpers.BcryptHelper) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		jwt: jwt,
		bcrypt: bcrypt,
	}
}

func (u *userUsecase) Register(req *dtos.CreateUserRequest) (*dtos.CreateUserResponse, error) {
	userExists, err := u.userRepo.GetByUsername(req.Username)
	if err == nil && userExists != nil {
		return nil, errors.New("username already taken")
	}

	role := model.Role(req.Role)
	if !role.IsValid() {
		return nil, errors.New("invalid role specified")
	}

	hashPassword, err := u.bcrypt.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &model.User{
		Username: req.Username,
		Password: hashPassword,
		FullName: req.FullName,
		Email: req.Email,
		Role: role,
		IsActive: true,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	res := &dtos.CreateUserResponse{
		ID: user.ID.String(),
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		Role: string(user.Role),
		IsActive: user.IsActive,
	}

	return res, nil
}

func (u *userUsecase) ListUsers(page, limit int) ([]*dtos.UserResponse, int64, error) {
	offset := (page - 1) * limit
	users, total, err := u.userRepo.FindAllPaginated(offset, limit)
	if err != nil {
		return nil, 0, errors.New("failed to fetch users")
	}

	var res []*dtos.UserResponse
	for _, user := range users {
		res = append(res, &dtos.UserResponse{
			ID:       user.ID.String(),
			Username: user.Username,
			FullName: user.FullName,
			Email:    user.Email,
			Role:     string(user.Role),
			IsActive: user.IsActive,
		})
	}

	return res, total, nil
}

func (u *userUsecase) Login(req *dtos.LoginRequest) (*dtos.LoginResponse, error) {
	user, err := u.userRepo.GetByUsername(req.Username)
	if err != nil || user == nil {
		return nil, errors.New("invalid username or password")
	}

	if !user.IsActive {
		return nil, errors.New("invalid username or password")
	}

	if err := u.bcrypt.ComparePassword(user.Password, req.Password); err != nil {
		return nil, errors.New("invalid username or password")
	}

	token, err := u.jwt.GenerateAccessToken(user.ID.String(), user.Username, string(user.Role), 1)
	if err != nil {
		return nil, errors.New("failed to generate JWT token")
	}

	res := &dtos.LoginResponse{
		AccessToken: token,
	}

	return res, nil
}
