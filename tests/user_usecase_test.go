package tests

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	dtos "github.com/savanyv/zenith-pay/internal/dto"
	"github.com/savanyv/zenith-pay/internal/model"
	"github.com/savanyv/zenith-pay/internal/usecase"
	"github.com/savanyv/zenith-pay/tests/mocks"
)

func TestUserUsecase_Register_Success(t *testing.T) {
	userRepo := &mocks.UserRepo{
		GetByUsernameFn: func(username string) (*model.User, error) {
			return nil, nil
		},
		CreateFn: func(user *model.User) error {
			user.ID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
			return nil
		},
	}
	bcrypt := &mocks.BcryptHelper{
		HashPasswordFn: func(password string) (string, error) {
			return "hashed", nil
		},
	}
	jwt := &mocks.JWTService{}
	uc := usecase.NewUserUsecase(userRepo, jwt, bcrypt)

	res, err := uc.Register(&dtos.CreateUserRequest{
		Username: "testuser",
		Password: "password",
		FullName: "Test User",
		Email:    "test@example.com",
		Role:     "cashier",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Username != "testuser" {
		t.Fatalf("expected testuser, got %s", res.Username)
	}
}

func TestUserUsecase_Register_UsernameTaken(t *testing.T) {
	userRepo := &mocks.UserRepo{
		GetByUsernameFn: func(username string) (*model.User, error) {
			return &model.User{Username: username}, nil
		},
	}
	uc := usecase.NewUserUsecase(userRepo, &mocks.JWTService{}, &mocks.BcryptHelper{})

	_, err := uc.Register(&dtos.CreateUserRequest{
		Username: "testuser",
		Password: "password",
		FullName: "Test User",
		Email:    "test@example.com",
		Role:     "cashier",
	})

	if err == nil || err.Error() != "username already taken" {
		t.Fatalf("expected username already taken, got %v", err)
	}
}

func TestUserUsecase_Register_InvalidRole(t *testing.T) {
	userRepo := &mocks.UserRepo{
		GetByUsernameFn: func(username string) (*model.User, error) {
			return nil, nil
		},
	}
	uc := usecase.NewUserUsecase(userRepo, &mocks.JWTService{}, &mocks.BcryptHelper{})

	_, err := uc.Register(&dtos.CreateUserRequest{
		Username: "testuser",
		Password: "password",
		FullName: "Test User",
		Email:    "test@example.com",
		Role:     "superadmin",
	})

	if err == nil || err.Error() != "invalid role specified" {
		t.Fatalf("expected invalid role specified, got %v", err)
	}
}

func TestUserUsecase_Register_FailHashPassword(t *testing.T) {
	userRepo := &mocks.UserRepo{
		GetByUsernameFn: func(username string) (*model.User, error) {
			return nil, nil
		},
	}
	bcrypt := &mocks.BcryptHelper{
		HashPasswordFn: func(password string) (string, error) {
			return "", errors.New("hash error")
		},
	}
	uc := usecase.NewUserUsecase(userRepo, &mocks.JWTService{}, bcrypt)

	_, err := uc.Register(&dtos.CreateUserRequest{
		Username: "testuser",
		Password: "password",
		FullName: "Test User",
		Email:    "test@example.com",
		Role:     "cashier",
	})

	if err == nil || err.Error() != "failed to hash password" {
		t.Fatalf("expected failed to hash password, got %v", err)
	}
}

func TestUserUsecase_Register_FailCreate(t *testing.T) {
	userRepo := &mocks.UserRepo{
		GetByUsernameFn: func(username string) (*model.User, error) {
			return nil, nil
		},
		CreateFn: func(user *model.User) error {
			return errors.New("db error")
		},
	}
	bcrypt := &mocks.BcryptHelper{
		HashPasswordFn: func(password string) (string, error) {
			return "hashed", nil
		},
	}
	uc := usecase.NewUserUsecase(userRepo, &mocks.JWTService{}, bcrypt)

	_, err := uc.Register(&dtos.CreateUserRequest{
		Username: "testuser",
		Password: "password",
		FullName: "Test User",
		Email:    "test@example.com",
		Role:     "cashier",
	})

	if err == nil || err.Error() != "failed to create user" {
		t.Fatalf("expected failed to create user, got %v", err)
	}
}

func TestUserUsecase_Login_Success(t *testing.T) {
	userRepo := &mocks.UserRepo{
		GetByUsernameFn: func(username string) (*model.User, error) {
			return &model.User{
				ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Username: "testuser",
				Password: "hashedpass",
				Role:     model.CashierRole,
				IsActive: true,
			}, nil
		},
	}
	bcrypt := &mocks.BcryptHelper{
		ComparePasswordFn: func(hashedPassword, password string) error {
			return nil
		},
	}
	jwt := &mocks.JWTService{
		GenerateAccessTokenFn: func(userID, username, role string, tokenVersion int) (string, error) {
			return "token123", nil
		},
	}
	uc := usecase.NewUserUsecase(userRepo, jwt, bcrypt)

	res, err := uc.Login(&dtos.LoginRequest{Username: "testuser", Password: "password"})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.AccessToken != "token123" {
		t.Fatalf("expected token123, got %s", res.AccessToken)
	}
}

func TestUserUsecase_Login_UserNotFound(t *testing.T) {
	userRepo := &mocks.UserRepo{
		GetByUsernameFn: func(username string) (*model.User, error) {
			return nil, errors.New("not found")
		},
	}
	uc := usecase.NewUserUsecase(userRepo, &mocks.JWTService{}, &mocks.BcryptHelper{})

	_, err := uc.Login(&dtos.LoginRequest{Username: "unknown", Password: "password"})

	if err == nil || err.Error() != "invalid username or password" {
		t.Fatalf("expected invalid username or password, got %v", err)
	}
}

func TestUserUsecase_Login_UserNotActive(t *testing.T) {
	userRepo := &mocks.UserRepo{
		GetByUsernameFn: func(username string) (*model.User, error) {
			return &model.User{
				IsActive: false,
			}, nil
		},
	}
	uc := usecase.NewUserUsecase(userRepo, &mocks.JWTService{}, &mocks.BcryptHelper{})

	_, err := uc.Login(&dtos.LoginRequest{Username: "testuser", Password: "password"})

	if err == nil || err.Error() != "invalid username or password" {
		t.Fatalf("expected invalid username or password, got %v", err)
	}
}

func TestUserUsecase_Login_PasswordMismatch(t *testing.T) {
	userRepo := &mocks.UserRepo{
		GetByUsernameFn: func(username string) (*model.User, error) {
			return &model.User{
				Password: "hashedpass",
				IsActive: true,
			}, nil
		},
	}
	bcrypt := &mocks.BcryptHelper{
		ComparePasswordFn: func(hashedPassword, password string) error {
			return errors.New("mismatch")
		},
	}
	uc := usecase.NewUserUsecase(userRepo, &mocks.JWTService{}, bcrypt)

	_, err := uc.Login(&dtos.LoginRequest{Username: "testuser", Password: "wrong"})

	if err == nil || err.Error() != "invalid username or password" {
		t.Fatalf("expected invalid username or password, got %v", err)
	}
}

func TestUserUsecase_Login_FailGenerateToken(t *testing.T) {
	userRepo := &mocks.UserRepo{
		GetByUsernameFn: func(username string) (*model.User, error) {
			return &model.User{
				ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Username: "testuser",
				Password: "hashedpass",
				Role:     model.CashierRole,
				IsActive: true,
			}, nil
		},
	}
	bcrypt := &mocks.BcryptHelper{
		ComparePasswordFn: func(hashedPassword, password string) error {
			return nil
		},
	}
	jwt := &mocks.JWTService{
		GenerateAccessTokenFn: func(userID, username, role string, tokenVersion int) (string, error) {
			return "", errors.New("token error")
		},
	}
	uc := usecase.NewUserUsecase(userRepo, jwt, bcrypt)

	_, err := uc.Login(&dtos.LoginRequest{Username: "testuser", Password: "password"})

	if err == nil || err.Error() != "failed to generate JWT token" {
		t.Fatalf("expected failed to generate JWT token, got %v", err)
	}
}
