package mocks

import "github.com/savanyv/zenith-pay/internal/model"

type ShiftRepo struct {
	CreateFn                 func(shift *model.Shift) error
	FindActiveShiftByCashierFn func(cashierID string) (*model.Shift, error)
	FindByIDFn               func(ID string) (*model.Shift, error)
	CloseShiftFn             func(shift *model.Shift) error
}

func (m *ShiftRepo) Create(shift *model.Shift) error                     { return m.CreateFn(shift) }
func (m *ShiftRepo) FindActiveShiftByCashier(cashierID string) (*model.Shift, error) {
	return m.FindActiveShiftByCashierFn(cashierID)
}
func (m *ShiftRepo) FindByID(ID string) (*model.Shift, error)            { return m.FindByIDFn(ID) }
func (m *ShiftRepo) CloseShift(shift *model.Shift) error                 { return m.CloseShiftFn(shift) }
