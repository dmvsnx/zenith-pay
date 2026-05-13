package mocks

type BcryptHelper struct {
	HashPasswordFn    func(password string) (string, error)
	ComparePasswordFn func(hashedPassword, password string) error
}

func (m *BcryptHelper) HashPassword(password string) (string, error) {
	return m.HashPasswordFn(password)
}

func (m *BcryptHelper) ComparePassword(hashedPassword, password string) error {
	return m.ComparePasswordFn(hashedPassword, password)
}
