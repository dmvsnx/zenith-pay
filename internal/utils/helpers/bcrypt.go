package helpers

import "golang.org/x/crypto/bcrypt"

type BcryptHelper interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
}

type bcryptHelper struct{
	bcryptCost int
}

func NewBcryptHelper() BcryptHelper {
	return &bcryptHelper{
		bcryptCost: bcrypt.DefaultCost,
	}
}

func (b *bcryptHelper) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), b.bcryptCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (b *bcryptHelper) ComparePassword(hashedPassword, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return err
	}

	return nil
}
